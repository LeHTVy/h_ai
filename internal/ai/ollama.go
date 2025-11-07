package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	baseURL    string
	model      string
	httpClient *http.Client
	logger     *zap.Logger
	enabled    bool
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model    string    `json:"model"`
	Prompt   string    `json:"prompt"`
	Stream   bool      `json:"stream,omitempty"`
	Options  *Options  `json:"options,omitempty"`
	Messages []Message `json:"messages,omitempty"`
}

// Options for Ollama API
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"` // Max tokens to generate
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Model              string   `json:"model"`
	CreatedAt          string   `json:"created_at"`
	Response           string   `json:"response"`
	Done               bool     `json:"done"`
	Context            []int    `json:"context,omitempty"`
	TotalDuration      int64    `json:"total_duration,omitempty"`
	LoadDuration       int64    `json:"load_duration,omitempty"`
	PromptEvalCount    int      `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64    `json:"prompt_eval_duration,omitempty"`
	EvalCount          int      `json:"eval_count,omitempty"`
	EvalDuration       int64    `json:"eval_duration,omitempty"`
	Message            *Message `json:"message,omitempty"`
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL string, model string, logger *zap.Logger) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama2" // Default model
	}

	client := &OllamaClient{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // 2 minutes timeout for LLM inference
		},
		logger:  logger,
		enabled: true,
	}

	// Test connection
	if err := client.Ping(); err != nil {
		logger.Warn("Ollama connection failed, AI features will be disabled",
			zap.String("error", err.Error()))
		client.enabled = false
	} else {
		logger.Info("Ollama client initialized",
			zap.String("baseURL", baseURL),
			zap.String("model", model))
	}

	return client
}

// Ping checks if Ollama server is reachable
func (c *OllamaClient) Ping() error {
	if !c.enabled {
		return fmt.Errorf("ollama client is disabled")
	}

	resp, err := c.httpClient.Get(c.baseURL + "/api/tags")
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama server returned status %d", resp.StatusCode)
	}

	return nil
}

// ModelInfo represents information about an Ollama model
type ModelInfo struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
}

// ModelsResponse represents the response from /api/tags
type ModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

// ListModels returns a list of available Ollama models
func (c *OllamaClient) ListModels() ([]string, error) {
	if !c.enabled {
		return nil, fmt.Errorf("ollama client is disabled")
	}

	resp, err := c.httpClient.Get(c.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var modelsResp ModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&modelsResp); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	models := make([]string, 0, len(modelsResp.Models))
	for _, model := range modelsResp.Models {
		models = append(models, model.Name)
	}

	return models, nil
}

// SetModel changes the model used by this client
func (c *OllamaClient) SetModel(model string) {
	c.model = model
	c.logger.Info("Ollama model changed", zap.String("model", model))
}

// GetModel returns the current model name
func (c *OllamaClient) GetModel() string {
	return c.model
}

// Generate generates text using Ollama
func (c *OllamaClient) Generate(prompt string, options *Options) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("ollama client is disabled")
	}

	req := OllamaRequest{
		Model:   c.model,
		Prompt:  prompt,
		Stream:  false,
		Options: options,
	}

	return c.generateInternal(req)
}

// Chat sends a chat message and gets a response
func (c *OllamaClient) Chat(messages []Message, options *Options) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("ollama client is disabled")
	}

	// Ensure we have default options with num_predict for complete responses
	if options == nil {
		options = &Options{
			Temperature: 0.7,
			TopP:        0.9,
			NumPredict:  4096, // Default max tokens for complete responses
		}
	} else if options.NumPredict == 0 {
		// If num_predict not set, use default
		options.NumPredict = 4096
	}

	// Add system message to encourage complete responses
	enhancedMessages := make([]Message, 0, len(messages)+1)

	// Check if first message is already a system message
	if len(messages) > 0 && messages[0].Role == "system" {
		enhancedMessages = append(enhancedMessages, messages[0])
		// Enhance existing system message
		enhancedMessages[0].Content = enhancedMessages[0].Content + "\n\nIMPORTANT: Always provide complete, detailed, and comprehensive answers. Do not cut off mid-sentence. Ensure your responses are fully formed and complete."
		enhancedMessages = append(enhancedMessages, messages[1:]...)
	} else {
		// Add new system message
		enhancedMessages = append(enhancedMessages, Message{
			Role:    "system",
			Content: "You are a helpful AI assistant. Always provide complete, detailed, and comprehensive answers. Do not cut off mid-sentence. Ensure your responses are fully formed and complete. Answer thoroughly and in full sentences.",
		})
		enhancedMessages = append(enhancedMessages, messages...)
	}

	req := OllamaRequest{
		Model:    c.model,
		Messages: enhancedMessages,
		Stream:   false,
		Options:  options,
	}

	return c.chatInternal(req)
}

// AnalyzeScanResults analyzes scan results and provides insights
func (c *OllamaClient) AnalyzeScanResults(tool string, results string, targetProfile string) (string, error) {
	if !c.enabled {
		return "", fmt.Errorf("ollama client is disabled")
	}

	prompt := fmt.Sprintf(`You are a cybersecurity expert. Analyze the following scan results and provide insights.

Tool Used: %s
Target: %s

Scan Results:
%s

Please provide:
1. Key findings and potential vulnerabilities
2. Recommended next steps
3. Risk assessment
4. Suggested follow-up tools or techniques

Be concise and actionable.`, tool, targetProfile, results)

	options := &Options{
		Temperature: 0.3, // Lower temperature for more focused analysis
		TopP:        0.9,
	}

	return c.Generate(prompt, options)
}

// SuggestTools suggests optimal tools based on target information
func (c *OllamaClient) SuggestTools(targetType string, technologies []string, objective string) ([]string, error) {
	if !c.enabled {
		return nil, fmt.Errorf("ollama client is disabled")
	}

	techStr := ""
	if len(technologies) > 0 {
		techStr = fmt.Sprintf("Technologies detected: %v. ", technologies)
	}

	prompt := fmt.Sprintf(`You are a penetration testing expert. Suggest the best security tools for the following scenario.

Target Type: %s
%sObjective: %s

Available tools: nmap, masscan, rustscan, gobuster, feroxbuster, ffuf, nuclei, nikto, sqlmap, wpscan, hydra, msfconsole, amass, subfinder, httpx

Provide a JSON array of 5-8 recommended tool names in order of priority. Return ONLY the JSON array, no additional text.
Example: ["nmap", "gobuster", "nuclei"]`, targetType, techStr, objective)

	options := &Options{
		Temperature: 0.4,
		TopP:        0.9,
	}

	response, err := c.Generate(prompt, options)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var tools []string
	if err := json.Unmarshal([]byte(response), &tools); err != nil {
		// If parsing fails, try to extract tool names manually
		c.logger.Warn("Failed to parse tool suggestions as JSON, using fallback",
			zap.String("response", response))
		return c.extractToolNames(response), nil
	}

	return tools, nil
}

// OptimizeParameters optimizes tool parameters using AI
func (c *OllamaClient) OptimizeParameters(tool string, target string, context map[string]interface{}) (map[string]interface{}, error) {
	if !c.enabled {
		return nil, fmt.Errorf("ollama client is disabled")
	}

	contextStr := ""
	for k, v := range context {
		contextStr += fmt.Sprintf("%s: %v\n", k, v)
	}

	prompt := fmt.Sprintf(`You are a security tool expert. Optimize parameters for %s.

Target: %s
Context:
%s

Provide optimized parameters as a JSON object. Return ONLY the JSON object, no additional text.
Example: {"target": "example.com", "scan_type": "-sV", "ports": "80,443"}

Ensure parameters are practical and effective.`, tool, target, contextStr)

	options := &Options{
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := c.Generate(prompt, options)
	if err != nil {
		return nil, err
	}

	// Parse JSON response
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(response), &params); err != nil {
		c.logger.Warn("Failed to parse parameters as JSON",
			zap.String("response", response),
			zap.Error(err))
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return params, nil
}

// IsEnabled returns whether Ollama client is enabled
func (c *OllamaClient) IsEnabled() bool {
	return c.enabled
}

// SetEnabled enables or disables the Ollama client
func (c *OllamaClient) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// internal methods

func (c *OllamaClient) generateInternal(req OllamaRequest) (string, error) {
	url := c.baseURL + "/api/generate"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read entire response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to decode as single JSON object
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		// If single object fails, might be streaming format
		return c.parseStreamingResponse(bodyBytes)
	}

	// Check if response is complete
	if ollamaResp.Done {
		return ollamaResp.Response, nil
	}

	// If not done, return what we have (should not happen with stream=false)
	return ollamaResp.Response, nil
}

// parseStreamingResponse parses newline-delimited JSON streaming response
func (c *OllamaClient) parseStreamingResponse(bodyBytes []byte) (string, error) {
	var fullResponse strings.Builder
	lines := strings.Split(string(bodyBytes), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var chunk OllamaResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			c.logger.Warn("Failed to parse streaming chunk", zap.Error(err), zap.String("chunk", line))
			continue
		}

		// Accumulate response
		if chunk.Message != nil {
			fullResponse.WriteString(chunk.Message.Content)
		} else if chunk.Response != "" {
			fullResponse.WriteString(chunk.Response)
		}

		// If done, we have the complete response
		if chunk.Done {
			break
		}
	}

	result := fullResponse.String()
	if result == "" {
		return "", fmt.Errorf("no valid response found in streaming data")
	}

	return result, nil
}

func (c *OllamaClient) chatInternal(req OllamaRequest) (string, error) {
	// Try /api/chat first (newer Ollama versions)
	url := c.baseURL + "/api/chat"

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// If /api/chat returns 405 (Method Not Allowed) or 404, fallback to /api/generate
	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotFound {
		c.logger.Debug("Ollama /api/chat not available, falling back to /api/generate",
			zap.Int("status", resp.StatusCode))

		// Convert chat messages to a prompt for /api/generate
		prompt := c.convertMessagesToPrompt(req.Messages)

		// Use generate instead
		generateReq := OllamaRequest{
			Model:   req.Model,
			Prompt:  prompt,
			Stream:  false,
			Options: req.Options,
		}

		return c.generateInternal(generateReq)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read entire response body first
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to decode as single JSON object first
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		// If single object fails, might be streaming format (newline-delimited JSON)
		// Try to parse as streaming response
		return c.parseStreamingResponse(bodyBytes)
	}

	// Check if response is complete
	if ollamaResp.Done {
		if ollamaResp.Message != nil {
			return ollamaResp.Message.Content, nil
		}
		return ollamaResp.Response, nil
	}

	// If not done, might be partial response - try to accumulate
	// This should not happen with stream=false, but handle it anyway
	if ollamaResp.Message != nil {
		return ollamaResp.Message.Content, nil
	}

	return ollamaResp.Response, nil
}

// convertMessagesToPrompt converts chat messages to a single prompt string
func (c *OllamaClient) convertMessagesToPrompt(messages []Message) string {
	var prompt strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case "system":
			prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
		case "user":
			prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
		case "assistant":
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
		}
	}

	prompt.WriteString("Assistant:")
	return prompt.String()
}

func (c *OllamaClient) extractToolNames(response string) []string {
	// Fallback: extract tool names from text
	tools := []string{}
	availableTools := []string{
		"nmap", "masscan", "rustscan", "gobuster", "feroxbuster",
		"ffuf", "nuclei", "nikto", "sqlmap", "wpscan", "hydra",
		"msfconsole", "amass", "subfinder", "httpx",
	}

	responseLower := fmt.Sprintf(" %s ", response)
	for _, tool := range availableTools {
		if contains(responseLower, " "+tool+" ") {
			tools = append(tools, tool)
		}
	}

	// Return top 8 tools found, or all if less than 8
	if len(tools) > 8 {
		return tools[:8]
	}
	return tools
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
