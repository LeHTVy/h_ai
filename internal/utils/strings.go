package utils

import (
	"strings"
)

// SanitizeCommand sanitizes command to prevent injection
func SanitizeCommand(cmd string) string {
	// Remove newlines and other potentially dangerous characters
	cmd = strings.ReplaceAll(cmd, "\n", "")
	cmd = strings.ReplaceAll(cmd, "\r", "")
	cmd = strings.TrimSpace(cmd)
	return cmd
}

// BuildCommand builds a command string from parts
func BuildCommand(parts ...string) string {
	return strings.Join(parts, " ")
}
