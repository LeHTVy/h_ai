package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LeHTVy/h_ai/internal/server"
	"go.uber.org/zap"
)

const (
	defaultPort = 8888
	defaultHost = "0.0.0.0"
)

func main() {
	var (
		port       = flag.Int("port", defaultPort, "Port for the API server")
		host       = flag.String("host", defaultHost, "Host for the API server")
		debug      = flag.Bool("debug", false, "Enable debug mode")
		ollamaURL  = flag.String("ollama-url", "", "Ollama API URL (default: http://localhost:11434)")
		ollamaModel = flag.String("ollama-model", "", "Ollama model to use (default: llama2)")
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

	// Print banner
	printBanner(*port, *debug, *ollamaURL, *ollamaModel)

	// Create and start server
	srv := server.New(*host, *port, logger, *ollamaURL, *ollamaModel)
	if err := srv.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func printBanner(port int, debug bool, ollamaURL string, ollamaModel string) {
	aiStatus := "Disabled"
	if ollamaURL != "" || ollamaModel != "" {
		aiStatus = fmt.Sprintf("Enabled (URL: %s, Model: %s)", 
			func() string { if ollamaURL == "" { return "default" }; return ollamaURL }(),
			func() string { if ollamaModel == "" { return "llama2" }; return ollamaModel }())
	}
	
	banner := fmt.Sprintf(`
██╗  ██╗███████╗██╗  ██╗███████╗████████╗██████╗ ██╗██╗  ██╗███████╗
██║  ██║██╔════╝╚██╗██╔╝██╔════╝╚══██╔══╝██╔══██╗██║██║ ██╔╝██╔════╝
███████║█████╗   ╚███╔╝ ███████╗   ██║   ██████╔╝██║█████╔╝ █████╗
██╔══██║██╔══╝   ██╔██╗ ╚════██║   ██║   ██╔══██╗██║██╔═██╗ ██╔══╝
██║  ██║███████╗██╔╝ ██╗███████║   ██║   ██║  ██║██║██║  ██╗███████╗
╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═╝╚══════╝

┌─────────────────────────────────────────────────────────────────────┐
│  🚀 H-AI - HexStrike AI Clone in Go                                 │
│  ⚡ AI-Automated Recon | Exploitation | Analysis Pipeline           │
│  🎯 Bug Bounty | CTF | Red Team | Zero-Day Research                │
└─────────────────────────────────────────────────────────────────────┘

[INFO] Server starting on %s:%d
[INFO] 150+ integrated modules | Adaptive AI decision engine active
[INFO] Ollama AI Integration: %s
[INFO] Golang implementation with enhanced performance
[INFO] Debug mode: %v
`, defaultHost, port, aiStatus, debug)
	fmt.Print(banner)
}
