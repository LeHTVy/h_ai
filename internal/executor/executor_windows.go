//go:build windows
// +build windows

package executor

import (
	"fmt"
	"os/exec"
)

// setUnixProcessGroup is a no-op on Windows
func setUnixProcessGroup(cmd *exec.Cmd) {
	// Windows doesn't support process groups the same way
}

// terminateProcess terminates a process on Windows
func terminateProcess(pid int) {
	// Windows: Use taskkill command
	killCmd := exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprintf("%d", pid))
	killCmd.Run() // Ignore errors, process might already be dead
}
