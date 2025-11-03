//go:build !windows
// +build !windows

package executor

import (
	"os/exec"
	"syscall"
	"time"
)

// setUnixProcessGroup sets process group for Unix systems
func setUnixProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Create process group
	}
}

// terminateProcess terminates a process on Unix systems
func terminateProcess(pid int) {
	// Try to kill process group first
	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
		time.Sleep(2 * time.Second)
		syscall.Kill(-pgid, syscall.SIGKILL)
	} else {
		// Fallback to direct kill
		syscall.Kill(pid, syscall.SIGTERM)
		time.Sleep(2 * time.Second)
		syscall.Kill(pid, syscall.SIGKILL)
	}
}
