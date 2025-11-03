package server

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/ai"
	"github.com/LeHTVy/h_ai/internal/cache"
	"github.com/LeHTVy/h_ai/internal/executor"
	"github.com/LeHTVy/h_ai/internal/intelligence"
	"github.com/LeHTVy/h_ai/internal/tools"
)

type Server struct {
	host     string
	port     int
	logger   *zap.Logger
	router   *gin.Engine
	httpSrv  *http.Server
	executor *executor.Executor
	cache    *cache.Cache
	tools    *tools.Manager
	engine   *intelligence.IntelligentDecisionEngine
}

func New(host string, port int, logger *zap.Logger, ollamaURL string, ollamaModel string) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(ginLogger(logger), gin.Recovery())

	cache := cache.New(30 * time.Minute, 10*time.Minute)
	exec := executor.New(logger, cache)
	toolsMgr := tools.New(logger, exec)
	
	// Initialize Ollama client (can be nil if not configured)
	var ollamaClient *ai.OllamaClient
	if ollamaURL != "" || ollamaModel != "" {
		ollamaClient = ai.NewOllamaClient(ollamaURL, ollamaModel, logger)
	}
	
	decisionEngine := intelligence.NewDecisionEngine(logger, ollamaClient)

	srv := &Server{
		host:     host,
		port:     port,
		logger:   logger,
		router:   router,
		executor: exec,
		cache:    cache,
		tools:    toolsMgr,
		engine:   decisionEngine,
	}

	srv.setupRoutes()
	return srv
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.handleHealth)

	// Chat web interface
	s.router.GET("/chat", s.handleChatPage)
	
	// Static files for chat UI
	s.router.Static("/static", "./web/static")

	// API routes
	api := s.router.Group("/api")
	{
		// Command execution
		api.POST("/command", s.handleCommand)

		// Tools endpoints
		tools := api.Group("/tools")
		{
			tools.POST("/nmap", s.handleNmap)
			tools.POST("/nmap-advanced", s.handleNmapAdvanced)
			tools.POST("/metasploit", s.handleMetasploit)
			tools.POST("/gobuster", s.handleGobuster)
			tools.POST("/nuclei", s.handleNuclei)
			tools.POST("/sqlmap", s.handleSqlmap)
			tools.POST("/hydra", s.handleHydra)
			tools.POST("/ffuf", s.handleFFuf)
			tools.POST("/netexec", s.handleNetexec)
			tools.POST("/amass", s.handleAmass)
			tools.POST("/masscan", s.handleMasscan)
			tools.POST("/autorecon", s.handleAutoRecon)
			tools.POST("/msfvenom", s.handleMSFVenom)
		}

		// Intelligence endpoints
		intel := api.Group("/intelligence")
		{
			intel.POST("/analyze-target", s.handleAnalyzeTarget)
			intel.POST("/select-tools", s.handleSelectTools)
			intel.POST("/optimize-parameters", s.handleOptimizeParameters)
			intel.POST("/create-attack-chain", s.handleCreateAttackChain)
			intel.POST("/smart-scan", s.handleSmartScan)
			intel.POST("/analyze-results", s.handleAnalyzeResults)
			intel.POST("/chat", s.handleAIChat)
		}

		// Process management
		process := api.Group("/processes")
		{
			process.GET("/list", s.handleProcessList)
			process.GET("/status/:pid", s.handleProcessStatus)
			process.POST("/terminate/:pid", s.handleProcessTerminate)
			process.GET("/dashboard", s.handleProcessDashboard)
		}

		// Cache endpoints
		cache := api.Group("/cache")
		{
			cache.GET("/stats", s.handleCacheStats)
		}

		// Telemetry
		api.GET("/telemetry", s.handleTelemetry)
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	s.logger.Info("Starting HTTP server", zap.String("address", addr))
	if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}

func ginLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		logger.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
		)
	}
}
