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
		port   = flag.Int("port", defaultPort, "Port for the API server")
		host   = flag.String("host", defaultHost, "Host for the API server")
		debug  = flag.Bool("debug", false, "Enable debug mode")
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
	printBanner(*port, *debug)

	// Create and start server
	srv := server.New(*host, *port, logger)
	if err := srv.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func printBanner(port int, debug bool) {
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
[INFO] Golang implementation with enhanced performance
[INFO] Debug mode: %v
`, defaultHost, port, debug)
	fmt.Print(banner)
}
