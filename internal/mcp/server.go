package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/client"
)

// MCP Request/Response structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError  `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]interface{} `json:"clientInfo"`
}

type InitializeResult struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      map[string]interface{} `json:"serverInfo"`
}

type Server struct {
	client     *client.Client
	logger     *zap.Logger
	reader     *bufio.Reader
	writer     *bufio.Writer
	initialized bool
	tools      []Tool
}

func NewServer(apiClient *client.Client, logger *zap.Logger) *Server {
	return &Server{
		client:     apiClient,
		logger:     logger,
		reader:     bufio.NewReader(os.Stdin),
		writer:     bufio.NewWriter(os.Stdout),
		initialized: false,
		tools:       buildTools(),
	}
}

func (s *Server) Run() error {
	s.logger.Info("Starting MCP server")

	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read error: %w", err)
		}

		line = trimNewline(line)
		if line == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.logger.Error("Failed to parse request", zap.Error(err))
			continue
		}

		resp := s.handleRequest(&req)
		if resp != nil {
			if err := s.sendResponse(resp); err != nil {
				s.logger.Error("Failed to send response", zap.Error(err))
			}
		}
	}

	return nil
}

func (s *Server) handleRequest(req *MCPRequest) *MCPResponse {
	s.logger.Debug("Handling request", zap.String("method", req.Method))

	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *Server) handleInitialize(req *MCPRequest) *MCPResponse {
	if s.initialized {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32000,
				Message: "Already initialized",
			},
		}
	}

	s.initialized = true

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		ServerInfo: map[string]interface{}{
			"name":    "h-ai",
			"version": "1.0.0",
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *Server) handleToolsList(req *MCPRequest) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": s.tools,
		},
	}
}

func (s *Server) handleToolsCall(req *MCPRequest) *MCPResponse {
	params, ok := req.Params.(map[string]interface{})
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Tool name required",
			},
		}
	}

	arguments, _ := params["arguments"].(map[string]interface{})

	result, err := s.executeTool(toolName, arguments)
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32000,
				Message: err.Error(),
			},
		}
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("%v", result),
				},
			},
		},
	}
}

func (s *Server) executeTool(toolName string, arguments map[string]interface{}) (interface{}, error) {
	s.logger.Info("Executing tool", zap.String("tool", toolName))

	switch toolName {
	case "nmap_scan":
		return s.executeNmap(arguments)
	case "gobuster_scan":
		return s.executeGobuster(arguments)
	case "metasploit_run":
		return s.executeMetasploit(arguments)
	case "nuclei_scan":
		return s.executeNuclei(arguments)
	case "sqlmap_scan":
		return s.executeSqlmap(arguments)
	case "hydra_attack":
		return s.executeHydra(arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", toolName)
	}
}

func (s *Server) executeNmap(arguments map[string]interface{}) (interface{}, error) {
	target, _ := arguments["target"].(string)
	if target == "" {
		return nil, fmt.Errorf("target is required")
	}

	scanType := getString(arguments, "scan_type", "-sV")
	ports := getString(arguments, "ports", "")
	additionalArgs := getString(arguments, "additional_args", "")

	data := map[string]interface{}{
		"target":         target,
		"scan_type":      scanType,
		"ports":          ports,
		"additional_args": additionalArgs,
		"use_recovery":   true,
	}

	result, err := s.client.Post("api/tools/nmap", data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) executeGobuster(arguments map[string]interface{}) (interface{}, error) {
	url, _ := arguments["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	mode := getString(arguments, "mode", "dir")
	wordlist := getString(arguments, "wordlist", "/usr/share/wordlists/dirb/common.txt")
	additionalArgs := getString(arguments, "additional_args", "")

	data := map[string]interface{}{
		"url":            url,
		"mode":           mode,
		"wordlist":       wordlist,
		"additional_args": additionalArgs,
	}

	result, err := s.client.Post("api/tools/gobuster", data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) executeMetasploit(arguments map[string]interface{}) (interface{}, error) {
	module, _ := arguments["module"].(string)
	if module == "" {
		return nil, fmt.Errorf("module is required")
	}

	options, _ := arguments["options"].(map[string]interface{})
	if options == nil {
		options = make(map[string]interface{})
	}

	data := map[string]interface{}{
		"module":  module,
		"options": options,
	}

	result, err := s.client.Post("api/tools/metasploit", data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) executeNuclei(arguments map[string]interface{}) (interface{}, error) {
	target, _ := arguments["target"].(string)
	if target == "" {
		return nil, fmt.Errorf("target is required")
	}

	data := map[string]interface{}{
		"target":          target,
		"templates":       getString(arguments, "templates", ""),
		"severity":        getString(arguments, "severity", ""),
		"additional_args": getString(arguments, "additional_args", ""),
	}

	result, err := s.client.Post("api/tools/nuclei", data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) executeSqlmap(arguments map[string]interface{}) (interface{}, error) {
	url, _ := arguments["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	data := map[string]interface{}{
		"url":            url,
		"data":           getString(arguments, "data", ""),
		"cookies":        getString(arguments, "cookies", ""),
		"additional_args": getString(arguments, "additional_args", ""),
	}

	result, err := s.client.Post("api/tools/sqlmap", data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) executeHydra(arguments map[string]interface{}) (interface{}, error) {
	target, _ := arguments["target"].(string)
	service, _ := arguments["service"].(string)
	if target == "" || service == "" {
		return nil, fmt.Errorf("target and service are required")
	}

	data := map[string]interface{}{
		"target":         target,
		"service":        service,
		"username":       getString(arguments, "username", ""),
		"password_list":  getString(arguments, "password_list", ""),
		"additional_args": getString(arguments, "additional_args", ""),
	}

	result, err := s.client.Post("api/tools/hydra", data)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Server) sendResponse(resp *MCPResponse) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if _, err := s.writer.Write(data); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	if _, err := s.writer.WriteString("\n"); err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return s.writer.Flush()
}

func buildTools() []Tool {
	return []Tool{
		{
			Name:        "nmap_scan",
			Description: "Execute an enhanced Nmap scan against a target",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"target":         map[string]interface{}{"type": "string", "description": "The IP address or hostname to scan"},
					"scan_type":      map[string]interface{}{"type": "string", "description": "Scan type (e.g., -sV for version detection)", "default": "-sV"},
					"ports":          map[string]interface{}{"type": "string", "description": "Comma-separated list of ports"},
					"additional_args": map[string]interface{}{"type": "string", "description": "Additional Nmap arguments"},
				},
				"required": []string{"target"},
			},
		},
		{
			Name:        "gobuster_scan",
			Description: "Execute Gobuster to find directories, DNS subdomains, or virtual hosts",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url":            map[string]interface{}{"type": "string", "description": "The target URL"},
					"mode":           map[string]interface{}{"type": "string", "description": "Scan mode (dir, dns, fuzz, vhost)", "default": "dir"},
					"wordlist":       map[string]interface{}{"type": "string", "description": "Path to wordlist file"},
					"additional_args": map[string]interface{}{"type": "string", "description": "Additional Gobuster arguments"},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "metasploit_run",
			Description: "Execute a Metasploit module",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"module":  map[string]interface{}{"type": "string", "description": "The Metasploit module to use"},
					"options": map[string]interface{}{"type": "object", "description": "Dictionary of module options"},
				},
				"required": []string{"module"},
			},
		},
		{
			Name:        "nuclei_scan",
			Description: "Execute Nuclei vulnerability scanner",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"target":         map[string]interface{}{"type": "string", "description": "Target URL or IP"},
					"templates":      map[string]interface{}{"type": "string", "description": "Template paths"},
					"severity":       map[string]interface{}{"type": "string", "description": "Severity level"},
					"additional_args": map[string]interface{}{"type": "string", "description": "Additional Nuclei arguments"},
				},
				"required": []string{"target"},
			},
		},
		{
			Name:        "sqlmap_scan",
			Description: "Execute SQLMap for SQL injection testing",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"url":            map[string]interface{}{"type": "string", "description": "Target URL"},
					"data":           map[string]interface{}{"type": "string", "description": "POST data"},
					"cookies":        map[string]interface{}{"type": "string", "description": "Cookies"},
					"additional_args": map[string]interface{}{"type": "string", "description": "Additional SQLMap arguments"},
				},
				"required": []string{"url"},
			},
		},
		{
			Name:        "hydra_attack",
			Description: "Execute Hydra for password brute forcing",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"target":         map[string]interface{}{"type": "string", "description": "Target IP or hostname"},
					"service":        map[string]interface{}{"type": "string", "description": "Service to attack (ssh, ftp, http, etc.)"},
					"username":       map[string]interface{}{"type": "string", "description": "Username"},
					"password_list":  map[string]interface{}{"type": "string", "description": "Password list file"},
					"additional_args": map[string]interface{}{"type": "string", "description": "Additional Hydra arguments"},
				},
				"required": []string{"target", "service"},
			},
		},
	}
}

func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		if val != "" {
			return val
		}
	}
	return defaultValue
}

func trimNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	if len(s) > 0 && s[len(s)-1] == '\r' {
		return s[:len(s)-1]
	}
	return s
}
