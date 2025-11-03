package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/client"
	"github.com/LeHTVy/h_ai/internal/mcp"
)

const (
	defaultServer = "http://127.0.0.1:8888"
	defaultTimeout = 300
)

func main() {
	var (
		server  = flag.String("server", defaultServer, "HexStrike AI API server URL")
		timeout = flag.Int("timeout", defaultTimeout, "Request timeout in seconds")
		debug   = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	// Initialize logger
	logConfig := zap.NewDevelopmentConfig()
	if !*debug {
		logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := logConfig.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting H-AI MCP Client",
		zap.String("server", *server),
		zap.Int("timeout", *timeout))

	// Create client
	apiClient := client.New(*server, *timeout, logger)

	// Check health
	health, err := apiClient.CheckHealth()
	if err != nil {
		logger.Warn("Failed to connect to API server", zap.Error(err))
		logger.Warn("MCP server will start, but tool execution may fail")
	} else {
		logger.Info("Successfully connected to API server",
			zap.Any("health", health))
	}

	// Start MCP server
	mcpServer := mcp.NewServer(apiClient, logger)
	
	logger.Info("MCP server ready to serve AI agents")
	
	if err := mcpServer.Run(); err != nil {
		logger.Fatal("MCP server failed", zap.Error(err))
	}
}
