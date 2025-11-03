package intelligence

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/ai"
)

// IntelligentDecisionEngine is AI-powered tool selection and parameter optimization engine
type IntelligentDecisionEngine struct {
	logger             *zap.Logger
	toolEffectiveness  map[string]map[string]float64
	technologySignatures map[string]map[string][]string
	attackPatterns     map[string][]AttackPattern
	ollamaClient      *ai.OllamaClient
}

// AttackPattern represents a pattern of tools to use
type AttackPattern struct {
	Tool     string                 `json:"tool"`
	Priority int                    `json:"priority"`
	Params   map[string]interface{} `json:"params"`
}

// NewDecisionEngine creates a new intelligent decision engine
func NewDecisionEngine(logger *zap.Logger, ollamaClient *ai.OllamaClient) *IntelligentDecisionEngine {
	engine := &IntelligentDecisionEngine{
		logger:        logger,
		ollamaClient: ollamaClient,
	}
	
	engine.toolEffectiveness = engine.initToolEffectiveness()
	engine.technologySignatures = engine.initTechnologySignatures()
	engine.attackPatterns = engine.initAttackPatterns()
	
	return engine
}

// AnalyzeTarget analyzes a target and creates a profile
func (e *IntelligentDecisionEngine) AnalyzeTarget(target string) *TargetProfile {
	profile := &TargetProfile{
		Target:           target,
		TargetType:       e.detectTargetType(target),
		IPAddresses:      []string{},
		OpenPorts:        []int{},
		Services:         make(map[int]string),
		Technologies:     []TechnologyStack{},
		SecurityHeaders:  make(map[string]string),
		SSLInfo:          make(map[string]interface{}),
		Subdomains:       []string{},
		Endpoints:        []string{},
		AttackSurfaceScore: 0.0,
		RiskLevel:        "unknown",
		ConfidenceScore:  0.5,
	}

	// Basic heuristics for target type detection
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		profile.TargetType = TargetTypeWebApplication
		profile.ConfidenceScore = 0.8
	} else if strings.Contains(target, ".") && !strings.Contains(target, "/") {
		// Likely an IP address or hostname
		profile.TargetType = TargetTypeNetworkHost
		profile.ConfidenceScore = 0.7
	}

	// Enhance analysis with AI if Ollama is available
	if e.ollamaClient != nil && e.ollamaClient.IsEnabled() {
		e.enhanceProfileWithAI(profile, target)
	}

	// Calculate attack surface score
	profile.AttackSurfaceScore = e.calculateAttackSurfaceScore(profile)
	
	// Determine risk level
	profile.RiskLevel = e.determineRiskLevel(profile.AttackSurfaceScore)

	return profile
}

// SelectOptimalTools selects optimal tools based on target profile and objective
func (e *IntelligentDecisionEngine) SelectOptimalTools(profile *TargetProfile, objective string) []string {
	var selectedTools []string

	// Use AI if available
	if e.ollamaClient != nil && e.ollamaClient.IsEnabled() {
		techList := make([]string, 0, len(profile.Technologies))
		for _, tech := range profile.Technologies {
			techList = append(techList, string(tech))
		}

		aiTools, err := e.ollamaClient.SuggestTools(string(profile.TargetType), techList, objective)
		if err == nil && len(aiTools) > 0 {
			e.logger.Info("AI suggested tools", zap.Strings("tools", aiTools))
			return aiTools
		}
		e.logger.Warn("AI tool suggestion failed, using rule-based fallback", zap.Error(err))
	}

	// Fallback to rule-based selection
	toolScores := make(map[string]float64)
	
	// Get effectiveness ratings for target type
	effectiveness, exists := e.toolEffectiveness[string(profile.TargetType)]
	if !exists {
		effectiveness = e.toolEffectiveness["unknown"]
	}

	// Score tools based on effectiveness
	for tool, score := range effectiveness {
		toolScores[tool] = score
	}

	// Adjust scores based on objective
	switch objective {
	case "quick":
		// Prioritize fast tools
		fastTools := []string{"nmap", "rustscan", "httpx", "gobuster"}
		for _, tool := range fastTools {
			if score, exists := toolScores[tool]; exists {
				toolScores[tool] = score * 1.2
			}
		}
	case "stealth":
		// Prioritize stealthy tools
		stealthTools := []string{"nmap-advanced", "masscan"}
		for _, tool := range stealthTools {
			if score, exists := toolScores[tool]; exists {
				toolScores[tool] = score * 1.15
			}
		}
	case "comprehensive":
		// Use all available tools
		// No adjustment needed
	}

	// Sort tools by score
	type toolScore struct {
		tool  string
		score float64
	}
	
	var sorted []toolScore
	for tool, score := range toolScores {
		sorted = append(sorted, toolScore{tool: tool, score: score})
	}
	
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})

	// Select top tools
	maxTools := 10
	if objective == "quick" {
		maxTools = 5
	} else if objective == "comprehensive" {
		maxTools = 15
	}

	selectedTools = make([]string, 0, maxTools)
	for i, ts := range sorted {
		if i >= maxTools {
			break
		}
		if ts.score > 0.5 { // Only select tools with effectiveness > 0.5
			selectedTools = append(selectedTools, ts.tool)
		}
	}

	return selectedTools
}

// OptimizeParameters optimizes tool parameters based on target profile
func (e *IntelligentDecisionEngine) OptimizeParameters(tool string, profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	// Try AI optimization first if available
	if e.ollamaClient != nil && e.ollamaClient.IsEnabled() {
		// Add profile info to context
		aiContext := make(map[string]interface{})
		for k, v := range context {
			aiContext[k] = v
		}
		aiContext["target_type"] = string(profile.TargetType)
		aiContext["open_ports"] = profile.OpenPorts
		
		aiParams, err := e.ollamaClient.OptimizeParameters(tool, profile.Target, aiContext)
		if err == nil && aiParams != nil && len(aiParams) > 0 {
			e.logger.Info("AI optimized parameters", 
				zap.String("tool", tool),
				zap.Any("params", aiParams))
			return aiParams
		}
		e.logger.Warn("AI parameter optimization failed, using rule-based fallback", zap.Error(err))
	}

	// Fallback to rule-based optimization
	params := make(map[string]interface{})

	switch tool {
	case "nmap":
		params = e.optimizeNmapParams(profile, context)
	case "gobuster":
		params = e.optimizeGobusterParams(profile, context)
	case "nuclei":
		params = e.optimizeNucleiParams(profile, context)
	case "sqlmap":
		params = e.optimizeSqlmapParams(profile, context)
	case "ffuf":
		params = e.optimizeFFufParams(profile, context)
	case "hydra":
		params = e.optimizeHydraParams(profile, context)
	default:
		// Default optimization
		params["target"] = profile.Target
	}

	return params
}

// CreateAttackChain creates an attack chain based on target profile and objective
func (e *IntelligentDecisionEngine) CreateAttackChain(profile *TargetProfile, objective string) *AttackChain {
	chain := &AttackChain{
		TargetProfile: profile,
		Steps:         []AttackStep{},
		RequiredTools: []string{},
		RiskLevel:     profile.RiskLevel,
	}

	// Get appropriate attack patterns
	patternKey := e.getPatternKey(profile.TargetType)
	patterns, exists := e.attackPatterns[patternKey]
	if !exists {
		patterns = e.attackPatterns["default"]
	}

	// Build attack chain from patterns
	for _, pattern := range patterns {
		step := AttackStep{
			Tool:                pattern.Tool,
			Parameters:          pattern.Params,
			ExpectedOutcome:     e.getExpectedOutcome(pattern.Tool),
			SuccessProbability:  e.getToolProbability(pattern.Tool, profile),
			ExecutionTimeEstimate: e.getEstimatedTime(pattern.Tool),
			Dependencies:        []string{},
		}
		
		// Add target to parameters
		step.Parameters["target"] = profile.Target
		
		chain.AddStep(step)
	}

	return chain
}

// Helper methods

func (e *IntelligentDecisionEngine) detectTargetType(target string) TargetType {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return TargetTypeWebApplication
	}
	if strings.Contains(target, "api") || strings.Contains(target, "/api/") {
		return TargetTypeAPIEndpoint
	}
	if strings.Contains(target, ".") && !strings.Contains(target, "/") {
		return TargetTypeNetworkHost
	}
	return TargetTypeUnknown
}

func (e *IntelligentDecisionEngine) calculateAttackSurfaceScore(profile *TargetProfile) float64 {
	score := 0.0
	
	// Port diversity
	score += float64(len(profile.OpenPorts)) * 0.1
	
	// Technology stack diversity
	score += float64(len(profile.Technologies)) * 0.15
	
	// Subdomain count
	score += float64(len(profile.Subdomains)) * 0.05
	
	// Endpoint count
	score += float64(len(profile.Endpoints)) * 0.08
	
	// Normalize score (0-1)
	return math.Min(score, 1.0)
}

func (e *IntelligentDecisionEngine) determineRiskLevel(score float64) string {
	if score >= 0.7 {
		return "high"
	} else if score >= 0.4 {
		return "medium"
	}
	return "low"
}

func (e *IntelligentDecisionEngine) getPatternKey(targetType TargetType) string {
	switch targetType {
	case TargetTypeWebApplication:
		return "web_reconnaissance"
	case TargetTypeAPIEndpoint:
		return "api_testing"
	case TargetTypeNetworkHost:
		return "network_discovery"
	default:
		return "default"
	}
}

func (e *IntelligentDecisionEngine) getExpectedOutcome(tool string) string {
	outcomes := map[string]string{
		"nmap":         "Port and service discovery",
		"gobuster":     "Directory and file enumeration",
		"nuclei":       "Vulnerability detection",
		"sqlmap":       "SQL injection vulnerabilities",
		"ffuf":         "Fuzzing results",
		"hydra":        "Credential discovery",
		"metasploit":   "Exploitation attempt",
		"masscan":      "Fast port scanning",
		"rustscan":     "Ultra-fast scanning",
	}
	
	if outcome, exists := outcomes[tool]; exists {
		return outcome
	}
	return "Tool execution results"
}

func (e *IntelligentDecisionEngine) getToolProbability(tool string, profile *TargetProfile) float64 {
	effectiveness, exists := e.toolEffectiveness[string(profile.TargetType)]
	if !exists {
		return 0.5
	}
	
	if prob, exists := effectiveness[tool]; exists {
		return prob
	}
	return 0.5
}

func (e *IntelligentDecisionEngine) getEstimatedTime(tool string) int {
	times := map[string]int{
		"nmap":         120,
		"gobuster":     300,
		"nuclei":       180,
		"sqlmap":       240,
		"ffuf":         200,
		"hydra":        600,
		"metasploit":   300,
		"masscan":      60,
		"rustscan":     30,
		"autorecon":    600,
	}
	
	if time, exists := times[tool]; exists {
		return time
	}
	return 120
}

// Parameter optimization methods for specific tools

func (e *IntelligentDecisionEngine) optimizeNmapParams(profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"target": profile.Target,
		"scan_type": "-sV",
	}
	
	if len(profile.OpenPorts) > 0 {
		// If we already know ports, scan those specifically
		ports := ""
		for i, port := range profile.OpenPorts {
			if i > 0 {
				ports += ","
			}
			ports += fmt.Sprintf("%d", port)
		}
		params["ports"] = ports
	} else {
		params["ports"] = "80,443,8080,8443,22,21,25,3306,5432"
	}
	
	params["additional_args"] = "-T4 -Pn"
	
	return params
}

func (e *IntelligentDecisionEngine) optimizeGobusterParams(profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"url": profile.Target,
		"mode": "dir",
		"wordlist": "/usr/share/wordlists/dirb/common.txt",
	}
	
	return params
}

func (e *IntelligentDecisionEngine) optimizeNucleiParams(profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"target": profile.Target,
		"severity": "critical,high",
	}
	
	return params
}

func (e *IntelligentDecisionEngine) optimizeSqlmapParams(profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"url": profile.Target,
	}
	
	return params
}

func (e *IntelligentDecisionEngine) optimizeFFufParams(profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"url": profile.Target,
		"wordlist": "/usr/share/wordlists/dirb/common.txt",
	}
	
	return params
}

func (e *IntelligentDecisionEngine) optimizeHydraParams(profile *TargetProfile, context map[string]interface{}) map[string]interface{} {
	params := map[string]interface{}{
		"target": profile.Target,
		"service": "http",
		"password_list": "/usr/share/wordlists/rockyou.txt",
	}
	
	return params
}

// Initialize tool effectiveness ratings
func (e *IntelligentDecisionEngine) initToolEffectiveness() map[string]map[string]float64 {
	return map[string]map[string]float64{
		string(TargetTypeWebApplication): {
			"nmap":     0.8,
			"gobuster": 0.9,
			"nuclei":   0.95,
			"sqlmap":   0.9,
			"ffuf":     0.9,
			"nikto":    0.85,
			"httpx":    0.85,
		},
		string(TargetTypeNetworkHost): {
			"nmap":         0.95,
			"nmap-advanced": 0.97,
			"masscan":      0.92,
			"rustscan":     0.9,
			"autorecon":    0.95,
			"hydra":        0.8,
			"netexec":      0.85,
		},
		string(TargetTypeAPIEndpoint): {
			"nuclei":     0.9,
			"ffuf":       0.85,
			"httpx":      0.9,
			"arjun":      0.95,
			"paramspider": 0.88,
		},
		"unknown": {
			"nmap": 0.8,
			"nuclei": 0.7,
		},
	}
}

// Initialize technology signatures
func (e *IntelligentDecisionEngine) initTechnologySignatures() map[string]map[string][]string {
	return map[string]map[string][]string{
		"headers": {
			"apache": {"Apache", "apache"},
			"nginx":  {"nginx", "Nginx"},
			"php":    {"PHP", "X-Powered-By: PHP"},
		},
		"content": {
			"wordpress": {"wp-content", "wp-includes", "WordPress"},
			"drupal":    {"Drupal", "drupal"},
		},
	}
}

// Initialize attack patterns
func (e *IntelligentDecisionEngine) initAttackPatterns() map[string][]AttackPattern {
	return map[string][]AttackPattern{
		"web_reconnaissance": {
			{Tool: "nmap", Priority: 1, Params: map[string]interface{}{"scan_type": "-sV", "ports": "80,443"}},
			{Tool: "gobuster", Priority: 2, Params: map[string]interface{}{"mode": "dir"}},
			{Tool: "nuclei", Priority: 3, Params: map[string]interface{}{"severity": "critical,high"}},
		},
		"api_testing": {
			{Tool: "httpx", Priority: 1, Params: map[string]interface{}{}},
			{Tool: "nuclei", Priority: 2, Params: map[string]interface{}{"tags": "api"}},
			{Tool: "ffuf", Priority: 3, Params: map[string]interface{}{}},
		},
		"network_discovery": {
			{Tool: "nmap", Priority: 1, Params: map[string]interface{}{"scan_type": "-sS"}},
			{Tool: "masscan", Priority: 2, Params: map[string]interface{}{}},
		},
		"default": {
			{Tool: "nmap", Priority: 1, Params: map[string]interface{}{}},
			{Tool: "nuclei", Priority: 2, Params: map[string]interface{}{}},
		},
	}
}

// enhanceProfileWithAI uses AI to enhance target profile analysis
func (e *IntelligentDecisionEngine) enhanceProfileWithAI(profile *TargetProfile, target string) {
	prompt := fmt.Sprintf(`You are a cybersecurity expert. Analyze this target and provide insights.

Target: %s
Detected Type: %s

Provide a JSON response with:
{
  "risk_assessment": "low/medium/high",
  "confidence": 0.0-1.0,
  "recommendations": ["recommendation1", "recommendation2"]
}

Return ONLY the JSON object, no additional text.`, target, profile.TargetType)

	options := &ai.Options{
		Temperature: 0.3,
		TopP:        0.9,
	}

	response, err := e.ollamaClient.Generate(prompt, options)
	if err != nil {
		e.logger.Warn("AI profile enhancement failed", zap.Error(err))
		return
	}

	// Try to parse JSON response (optional - non-critical)
	// If parsing fails, we continue with rule-based analysis
	e.logger.Debug("AI enhanced profile analysis", zap.String("response", response))
}

// AnalyzeScanResults analyzes scan results using AI
func (e *IntelligentDecisionEngine) AnalyzeScanResults(tool string, results string, target string) (string, error) {
	if e.ollamaClient == nil || !e.ollamaClient.IsEnabled() {
		return "", fmt.Errorf("Ollama AI is not available")
	}

	profile := e.AnalyzeTarget(target)
	profileStr := fmt.Sprintf("Target Type: %s, Risk Level: %s", profile.TargetType, profile.RiskLevel)
	
	return e.ollamaClient.AnalyzeScanResults(tool, results, profileStr)
}

// GetOllamaClient returns the Ollama client if available
func (e *IntelligentDecisionEngine) GetOllamaClient() *ai.OllamaClient {
	return e.ollamaClient
}
