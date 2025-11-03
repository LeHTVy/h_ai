package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/cache"
)

type ExecutionResult struct {
	Success      bool          `json:"success"`
	Stdout       string        `json:"stdout"`
	Stderr       string        `json:"stderr"`
	ReturnCode   int           `json:"return_code"`
	ExecutionTime float64     `json:"execution_time"`
	PID          int           `json:"pid,omitempty"`
}

type ProcessInfo struct {
	PID         int       `json:"pid"`
	Command     string    `json:"command"`
	StartTime   time.Time `json:"start_time"`
	Status      string    `json:"status"`
}

type Executor struct {
	logger      *zap.Logger
	cache       *cache.Cache
	processes   map[int]*ProcessInfo
	processLock sync.RWMutex
	timeout     time.Duration
}

func New(logger *zap.Logger, cache *cache.Cache) *Executor {
	executor := &Executor{
		logger:    logger,
		cache:     cache,
		processes: make(map[int]*ProcessInfo),
		timeout:   300 * time.Second, // 5 minutes default
	}

	// Handle cleanup on exit
	go executor.handleCleanup()
	return executor
}

func (e *Executor) Execute(command string, useCache bool) ExecutionResult {
	// Check cache first
	if useCache {
		if cached, found := e.cache.Get(command); found {
			if result, ok := cached.(ExecutionResult); ok {
				e.logger.Debug("Using cached result", zap.String("command", command))
				return result
			}
		}
	}

	start := time.Now()
	result := e.executeCommand(command)
	executionTime := time.Since(start).Seconds()
	result.ExecutionTime = executionTime

	// Cache successful results
	if useCache && result.Success {
		e.cache.Set(command, result, 30*time.Minute)
	}

	return result
}

func (e *Executor) executeCommand(command string) ExecutionResult {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Create process group
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return ExecutionResult{
			Success:    false,
			Stderr:     err.Error(),
			ReturnCode: -1,
		}
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return ExecutionResult{
			Success:    false,
			Stderr:     err.Error(),
			ReturnCode: -1,
		}
	}

	if err := cmd.Start(); err != nil {
		return ExecutionResult{
			Success:    false,
			Stderr:     err.Error(),
			ReturnCode: -1,
		}
	}

	pid := cmd.Process.Pid
	e.registerProcess(pid, command)

	// Read output in parallel
	var stdoutBytes, stderrBytes []byte
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				stdoutBytes = append(stdoutBytes, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		buf := make([]byte, 4096)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				stderrBytes = append(stderrBytes, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
	}()

	err = cmd.Wait()
	wg.Wait()

	e.unregisterProcess(pid)

	returnCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			returnCode = exitError.ExitCode()
		}
	}

	return ExecutionResult{
		Success:    returnCode == 0,
		Stdout:     string(stdoutBytes),
		Stderr:     string(stderrBytes),
		ReturnCode: returnCode,
		PID:        pid,
	}
}

func (e *Executor) registerProcess(pid int, command string) {
	e.processLock.Lock()
	defer e.processLock.Unlock()

	e.processes[pid] = &ProcessInfo{
		PID:       pid,
		Command:   command,
		StartTime: time.Now(),
		Status:    "running",
	}
}

func (e *Executor) unregisterProcess(pid int) {
	e.processLock.Lock()
	defer e.processLock.Unlock()

	if proc, exists := e.processes[pid]; exists {
		proc.Status = "completed"
		delete(e.processes, pid)
	}
}

func (e *Executor) ListProcesses() []ProcessInfo {
	e.processLock.RLock()
	defer e.processLock.RUnlock()

	processes := make([]ProcessInfo, 0, len(e.processes))
	for _, proc := range e.processes {
		processes = append(processes, *proc)
	}
	return processes
}

func (e *Executor) GetProcessStatus(pid int) *ProcessInfo {
	e.processLock.RLock()
	defer e.processLock.RUnlock()

	if proc, exists := e.processes[pid]; exists {
		return proc
	}
	return nil
}

func (e *Executor) TerminateProcess(pid int) error {
	e.processLock.RLock()
	proc, exists := e.processes[pid]
	e.processLock.RUnlock()

	if !exists {
		return fmt.Errorf("process %d not found", pid)
	}

	// Try to kill the process group
	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		syscall.Kill(-pgid, syscall.SIGTERM)
		time.Sleep(2 * time.Second)
		syscall.Kill(-pgid, syscall.SIGKILL)
	} else {
		syscall.Kill(pid, syscall.SIGTERM)
		time.Sleep(2 * time.Second)
		syscall.Kill(pid, syscall.SIGKILL)
	}

	e.unregisterProcess(pid)
	return nil
}

func (e *Executor) GetDashboard() map[string]interface{} {
	processes := e.ListProcesses()
	return map[string]interface{}{
		"active_processes": len(processes),
		"processes":        processes,
		"timestamp":        time.Now().Unix(),
	}
}

func (e *Executor) handleCleanup() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	e.logger.Info("Shutting down executor, terminating all processes")

	e.processLock.Lock()
	defer e.processLock.Unlock()

	for pid := range e.processes {
		e.TerminateProcess(pid)
	}
}
