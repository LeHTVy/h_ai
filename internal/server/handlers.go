package server

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/ai"
	"github.com/LeHTVy/h_ai/internal/models"
)

func (s *Server) handleHealth(c *gin.Context) {
	health := map[string]interface{}{
		"status":                  "healthy",
		"version":                 "1.0.0",
		"timestamp":               "",
		"tools_status":            s.tools.CheckToolsAvailability(),
		"all_essential_tools_available": true,
	}
	c.JSON(http.StatusOK, health)
}

func (s *Server) handleCommand(c *gin.Context) {
	var req models.CommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.executor.Execute(req.Command, req.UseCache)
	c.JSON(http.StatusOK, result)
}

// Nmap handler
func (s *Server) handleNmap(c *gin.Context) {
	var req models.NmapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target parameter is required"})
		return
	}

	result := s.tools.ExecuteNmap(req)
	c.JSON(http.StatusOK, result)
}

// Nmap Advanced handler
func (s *Server) handleNmapAdvanced(c *gin.Context) {
	var req models.NmapAdvancedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Target == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Target parameter is required"})
		return
	}

	result := s.tools.ExecuteNmapAdvanced(req)
	c.JSON(http.StatusOK, result)
}

// Metasploit handler
func (s *Server) handleMetasploit(c *gin.Context) {
	var req models.MetasploitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Module == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module parameter is required"})
		return
	}

	result := s.tools.ExecuteMetasploit(req)
	c.JSON(http.StatusOK, result)
}

// Gobuster handler
func (s *Server) handleGobuster(c *gin.Context) {
	var req models.GobusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteGobuster(req)
	c.JSON(http.StatusOK, result)
}

// Nuclei handler
func (s *Server) handleNuclei(c *gin.Context) {
	var req models.NucleiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteNuclei(req)
	c.JSON(http.StatusOK, result)
}

// SQLMap handler
func (s *Server) handleSqlmap(c *gin.Context) {
	var req models.SqlmapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteSqlmap(req)
	c.JSON(http.StatusOK, result)
}

// Hydra handler
func (s *Server) handleHydra(c *gin.Context) {
	var req models.HydraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteHydra(req)
	c.JSON(http.StatusOK, result)
}

// FFuf handler
func (s *Server) handleFFuf(c *gin.Context) {
	var req models.FFufRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteFFuf(req)
	c.JSON(http.StatusOK, result)
}

// NetExec handler
func (s *Server) handleNetexec(c *gin.Context) {
	var req models.NetexecRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteNetexec(req)
	c.JSON(http.StatusOK, result)
}

// Amass handler
func (s *Server) handleAmass(c *gin.Context) {
	var req models.AmassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteAmass(req)
	c.JSON(http.StatusOK, result)
}

// Masscan handler
func (s *Server) handleMasscan(c *gin.Context) {
	var req models.MasscanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteMasscan(req)
	c.JSON(http.StatusOK, result)
}

// AutoRecon handler
func (s *Server) handleAutoRecon(c *gin.Context) {
	var req models.AutoReconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteAutoRecon(req)
	c.JSON(http.StatusOK, result)
}

// MSFVenom handler
func (s *Server) handleMSFVenom(c *gin.Context) {
	var req models.MSFVenomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := s.tools.ExecuteMSFVenom(req)
	c.JSON(http.StatusOK, result)
}

// Intelligence handlers
func (s *Server) handleAnalyzeTarget(c *gin.Context) {
	var req models.AnalyzeTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Analyzing target", zap.String("target", req.Target))
	
	profile := s.engine.AnalyzeTarget(req.Target)
	
	result := map[string]interface{}{
		"success":       true,
		"target":        req.Target,
		"target_profile": profile,
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleSelectTools(c *gin.Context) {
	var req models.SelectToolsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Selecting tools", zap.String("target", req.Target))
	
	// Analyze target first
	profile := s.engine.AnalyzeTarget(req.Target)
	
	// Determine objective
	objective := req.TargetType
	if objective == "" {
		objective = "comprehensive"
	}
	
	// Select optimal tools
	selectedTools := s.engine.SelectOptimalTools(profile, objective)
	
	result := map[string]interface{}{
		"success":        true,
		"target":         req.Target,
		"target_profile": profile,
		"selected_tools": selectedTools,
		"tool_count":     len(selectedTools),
		"objective":      objective,
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleOptimizeParameters(c *gin.Context) {
	var req models.OptimizeParametersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Optimizing parameters", zap.String("tool", req.Tool))
	
	// Create a basic profile from context
	profile := s.engine.AnalyzeTarget("")
	if target, ok := req.Context["target"].(string); ok {
		profile = s.engine.AnalyzeTarget(target)
	}
	
	// Optimize parameters
	optimizedParams := s.engine.OptimizeParameters(req.Tool, profile, req.Context)
	
	result := map[string]interface{}{
		"success":         true,
		"tool":            req.Tool,
		"optimized_params": optimizedParams,
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleCreateAttackChain(c *gin.Context) {
	var req models.AnalyzeTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Creating attack chain", zap.String("target", req.Target))
	
	objective := req.AnalysisType
	if objective == "" {
		objective = "comprehensive"
	}
	
	profile := s.engine.AnalyzeTarget(req.Target)
	chain := s.engine.CreateAttackChain(profile, objective)
	
	result := map[string]interface{}{
		"success":      true,
		"target":       req.Target,
		"attack_chain": chain,
		"objective":    objective,
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleSmartScan(c *gin.Context) {
	var req models.AnalyzeTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Starting smart scan", zap.String("target", req.Target))
	
	objective := req.AnalysisType
	if objective == "" {
		objective = "comprehensive"
	}
	
	// Analyze target
	profile := s.engine.AnalyzeTarget(req.Target)
	
	// Select optimal tools
	selectedTools := s.engine.SelectOptimalTools(profile, objective)
	
	// Limit tools if needed
	maxTools := 5
	if len(selectedTools) > maxTools {
		selectedTools = selectedTools[:maxTools]
	}
	
	result := map[string]interface{}{
		"success":        true,
		"target":         req.Target,
		"target_profile": profile,
		"selected_tools":  selectedTools,
		"tool_count":     len(selectedTools),
		"objective":      objective,
		"message":        "Smart scan prepared. Tools are ready for execution.",
	}
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleAnalyzeResults(c *gin.Context) {
	var req struct {
		Tool     string `json:"tool" binding:"required"`
		Results  string `json:"results" binding:"required"`
		Target   string `json:"target" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.logger.Info("Analyzing scan results with AI", 
		zap.String("tool", req.Tool),
		zap.String("target", req.Target))

	// Use engine to analyze results with AI
	analysis, err := s.engine.AnalyzeScanResults(req.Tool, req.Results, req.Target)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "AI analysis requires Ollama integration. Please configure --ollama-url and --ollama-model flags when starting the server.",
		})
		return
	}

	result := map[string]interface{}{
		"success":  true,
		"tool":     req.Tool,
		"target":   req.Target,
		"analysis": analysis,
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleAIChat(c *gin.Context) {
	var req struct {
		Messages []struct {
			Role    string `json:"role" binding:"required"`
			Content string `json:"content" binding:"required"`
		} `json:"messages" binding:"required"`
		Temperature *float64 `json:"temperature,omitempty"`
		TopP        *float64 `json:"top_p,omitempty"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Ollama client from engine
	ollamaClient := s.engine.GetOllamaClient()
	if ollamaClient == nil || !ollamaClient.IsEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Ollama AI is not available",
			"message": "Please configure --ollama-url and --ollama-model flags when starting the server.",
		})
		return
	}

	// Convert request messages to AI messages
	aiMessages := make([]ai.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		aiMessages = append(aiMessages, ai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Set options
	options := &ai.Options{}
	if req.Temperature != nil {
		options.Temperature = *req.Temperature
	} else {
		options.Temperature = 0.7 // Default temperature
	}
	if req.TopP != nil {
		options.TopP = *req.TopP
	} else {
		options.TopP = 0.9 // Default top_p
	}

	s.logger.Info("Processing AI chat request", 
		zap.Int("message_count", len(aiMessages)))

	// Call Ollama chat
	response, err := ollamaClient.Chat(aiMessages, options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	result := map[string]interface{}{
		"success": true,
		"response": response,
		"message_count": len(aiMessages),
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) handleChatPage(c *gin.Context) {
	// Serve the chat HTML page
	chatHTMLPath := filepath.Join(".", "web", "static", "chat.html")
	c.File(chatHTMLPath)
}

// Process management handlers
func (s *Server) handleProcessList(c *gin.Context) {
	processes := s.executor.ListProcesses()
	c.JSON(http.StatusOK, gin.H{"processes": processes})
}

func (s *Server) handleProcessStatus(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PID"})
		return
	}

	status := s.executor.GetProcessStatus(pid)
	if status == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process not found"})
		return
	}

	c.JSON(http.StatusOK, status)
}

func (s *Server) handleProcessTerminate(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PID"})
		return
	}

	if err := s.executor.TerminateProcess(pid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process terminated", "pid": pid})
}

func (s *Server) handleProcessDashboard(c *gin.Context) {
	dashboard := s.executor.GetDashboard()
	c.JSON(http.StatusOK, dashboard)
}

// Cache handlers
func (s *Server) handleCacheStats(c *gin.Context) {
	stats := s.cache.Stats()
	c.JSON(http.StatusOK, stats)
}

// Telemetry handler
func (s *Server) handleTelemetry(c *gin.Context) {
	telemetry := map[string]interface{}{
		"uptime":           time.Since(time.Now()).Seconds(), // TODO: Track actual uptime
		"cpu_usage":        0,                                 // TODO: Implement CPU usage tracking
		"memory_usage":     0,                                 // TODO: Implement memory usage tracking
		"active_processes": len(s.executor.ListProcesses()),
	}
	c.JSON(http.StatusOK, telemetry)
}
