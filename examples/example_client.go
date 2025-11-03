package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/LeHTVy/h_ai/internal/client"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create client
	apiClient := client.New("http://127.0.0.1:8888", 300, logger)

	// Check health
	health, err := apiClient.CheckHealth()
	if err != nil {
		log.Fatal("Failed to check health:", err)
	}
	fmt.Println("Server health:", prettyPrint(health))

	// Execute Nmap scan
	nmapResult, err := apiClient.Post("api/tools/nmap", map[string]interface{}{
		"target":        "example.com",
		"scan_type":     "-sV",
		"ports":         "80,443",
		"additional_args": "-T4",
	})
	if err != nil {
		log.Fatal("Failed to execute nmap:", err)
	}
	fmt.Println("Nmap result:", prettyPrint(nmapResult))

	// Execute Gobuster scan
	gobusterResult, err := apiClient.Post("api/tools/gobuster", map[string]interface{}{
		"url":   "https://example.com",
		"mode":  "dir",
		"wordlist": "/usr/share/wordlists/dirb/common.txt",
	})
	if err != nil {
		log.Fatal("Failed to execute gobuster:", err)
	}
	fmt.Println("Gobuster result:", prettyPrint(gobusterResult))
}

func prettyPrint(data interface{}) string {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", data)
	}
	return string(jsonData)
}
