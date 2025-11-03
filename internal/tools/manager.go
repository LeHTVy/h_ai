package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/LeHTVy/h_ai/internal/executor"
	"github.com/LeHTVy/h_ai/internal/models"
)

type Manager struct {
	logger      *zap.Logger
	executor    *executor.Executor
	toolCache   map[string]bool
	cacheLock   sync.RWMutex
	toolTimeout int // in seconds
}

func New(logger *zap.Logger, exec *executor.Executor) *Manager {
	mgr := &Manager{
		logger:      logger,
		executor:    exec,
		toolCache:   make(map[string]bool),
		toolTimeout: 300,
	}

	// Pre-check tool availability
	go mgr.checkToolsAvailability()
	return mgr
}

func (m *Manager) CheckToolsAvailability() map[string]bool {
	m.cacheLock.RLock()
	defer m.cacheLock.RUnlock()

	// Return a copy
	result := make(map[string]bool)
	for tool, available := range m.toolCache {
		result[tool] = available
	}
	return result
}

func (m *Manager) checkToolsAvailability() {
	tools := []string{
		"nmap", "masscan", "rustscan", "gobuster", "feroxbuster",
		"ffuf", "nuclei", "nikto", "sqlmap", "wpscan", "hydra",
		"msfconsole", "msfvenom", "nxc", "amass", "subfinder",
	}

	m.cacheLock.Lock()
	defer m.cacheLock.Unlock()

	for _, tool := range tools {
		available := m.isToolAvailable(tool)
		m.toolCache[tool] = available
		if available {
			m.logger.Debug("Tool available", zap.String("tool", tool))
		} else {
			m.logger.Warn("Tool not available", zap.String("tool", tool))
		}
	}
}

func (m *Manager) isToolAvailable(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

func (m *Manager) buildCommand(tool string, args ...string) string {
	cmd := tool
	for _, arg := range args {
		cmd += " " + arg
	}
	return cmd
}

// ExecuteNmap executes an Nmap scan
func (m *Manager) ExecuteNmap(req models.NmapRequest) map[string]interface{} {
	scanType := req.ScanType
	if scanType == "" {
		scanType = "-sCV"
	}

	args := []string{scanType}
	if req.Ports != "" {
		args = append(args, "-p", req.Ports)
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}
	args = append(args, req.Target)

	command := m.buildCommand("nmap", args...)
	m.logger.Info("Executing Nmap scan", zap.String("target", req.Target))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteNmapAdvanced executes an advanced Nmap scan
func (m *Manager) ExecuteNmapAdvanced(req models.NmapAdvancedRequest) map[string]interface{} {
	scanType := req.ScanType
	if scanType == "" {
		scanType = "-sS"
	}

	args := []string{scanType, req.Target}
	if req.Ports != "" {
		args = append(args, "-p", req.Ports)
	}
	if req.Stealth {
		args = append(args, "-T2", "-f", "--mtu", "24")
	} else {
		timing := req.Timing
		if timing == "" {
			timing = "T4"
		}
		args = append(args, "-"+timing)
	}
	if req.OSDetection {
		args = append(args, "-O")
	}
	if req.VersionDetection {
		args = append(args, "-sV")
	}
	if req.Aggressive {
		args = append(args, "-A")
	}
	if req.NSEScripts != "" {
		args = append(args, "--script="+req.NSEScripts)
	} else if !req.Aggressive {
		args = append(args, "--script=default,discovery,safe")
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("nmap", args...)
	m.logger.Info("Executing Advanced Nmap scan", zap.String("target", req.Target))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteMetasploit executes a Metasploit module
func (m *Manager) ExecuteMetasploit(req models.MetasploitRequest) map[string]interface{} {
	// Create resource script
	resourceContent := fmt.Sprintf("use %s\n", req.Module)
	for key, value := range req.Options {
		resourceContent += fmt.Sprintf("set %s %s\n", key, value)
	}
	resourceContent += "exploit\n"

	// Save to temporary file
	resourceFile := filepath.Join(os.TempDir(), "h_ai_msf_resource.rc")
	if err := os.WriteFile(resourceFile, []byte(resourceContent), 0644); err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("Failed to create resource file: %v", err),
		}
	}
	defer os.Remove(resourceFile)

	command := fmt.Sprintf("msfconsole -q -r %s", resourceFile)
	m.logger.Info("Executing Metasploit module", zap.String("module", req.Module))

	result := m.executor.Execute(command, false)
	return m.formatResult(result)
}

// ExecuteGobuster executes a Gobuster scan
func (m *Manager) ExecuteGobuster(req models.GobusterRequest) map[string]interface{} {
	mode := req.Mode
	if mode == "" {
		mode = "dir"
	}

	wordlist := req.Wordlist
	if wordlist == "" {
		wordlist = "/usr/share/wordlists/dirb/common.txt"
	}

	args := []string{mode, "-u", req.URL, "-w", wordlist}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("gobuster", args...)
	m.logger.Info("Executing Gobuster scan", zap.String("url", req.URL))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteNuclei executes a Nuclei scan
func (m *Manager) ExecuteNuclei(req models.NucleiRequest) map[string]interface{} {
	args := []string{"-u", req.Target}
	if req.Templates != "" {
		args = append(args, "-t", req.Templates)
	}
	if req.Severity != "" {
		args = append(args, "-severity", req.Severity)
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("nuclei", args...)
	m.logger.Info("Executing Nuclei scan", zap.String("target", req.Target))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteSqlmap executes a SQLMap scan
func (m *Manager) ExecuteSqlmap(req models.SqlmapRequest) map[string]interface{} {
	args := []string{"-u", req.URL}
	if req.Data != "" {
		args = append(args, "--data", req.Data)
	}
	if req.Cookies != "" {
		args = append(args, "--cookie", req.Cookies)
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("sqlmap", args...)
	m.logger.Info("Executing SQLMap scan", zap.String("url", req.URL))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteHydra executes a Hydra brute force attack
func (m *Manager) ExecuteHydra(req models.HydraRequest) map[string]interface{} {
	args := []string{}
	if req.Username != "" {
		args = append(args, "-l", req.Username)
	} else {
		args = append(args, "-L", req.PasswordList) // Username list
	}
	if req.PasswordList != "" {
		args = append(args, "-P", req.PasswordList)
	}
	args = append(args, req.Target, req.Service)
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("hydra", args...)
	m.logger.Info("Executing Hydra attack", zap.String("target", req.Target))

	result := m.executor.Execute(command, false) // Don't cache brute force results
	return m.formatResult(result)
}

// ExecuteFFuf executes an FFuf fuzzing scan
func (m *Manager) ExecuteFFuf(req models.FFufRequest) map[string]interface{} {
	wordlist := req.Wordlist
	if wordlist == "" {
		wordlist = "/usr/share/wordlists/dirb/common.txt"
	}

	args := []string{"-u", req.URL + "/FUZZ", "-w", wordlist}
	if len(req.Headers) > 0 {
		for key, value := range req.Headers {
			args = append(args, "-H", fmt.Sprintf("%s: %s", key, value))
		}
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("ffuf", args...)
	m.logger.Info("Executing FFuf scan", zap.String("url", req.URL))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteNetexec executes a NetExec scan
func (m *Manager) ExecuteNetexec(req models.NetexecRequest) map[string]interface{} {
	protocol := req.Protocol
	if protocol == "" {
		protocol = "smb"
	}

	args := []string{protocol, req.Target}
	if req.Username != "" {
		args = append(args, "-u", req.Username)
	}
	if req.Password != "" {
		args = append(args, "-p", req.Password)
	}
	if req.Hash != "" {
		args = append(args, "-H", req.Hash)
	}
	if req.Module != "" {
		args = append(args, "-M", req.Module)
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("nxc", args...)
	m.logger.Info("Executing NetExec scan", zap.String("target", req.Target))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteAmass executes an Amass enumeration
func (m *Manager) ExecuteAmass(req models.AmassRequest) map[string]interface{} {
	args := []string{"enum", "-d", req.Domain}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("amass", args...)
	m.logger.Info("Executing Amass enumeration", zap.String("domain", req.Domain))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteMasscan executes a Masscan scan
func (m *Manager) ExecuteMasscan(req models.MasscanRequest) map[string]interface{} {
	ports := req.Ports
	if ports == "" {
		ports = "1-65535"
	}

	rate := req.Rate
	if rate == "" {
		rate = "1000"
	}

	args := []string{"-p", ports, "--rate", rate, req.Target}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("masscan", args...)
	m.logger.Info("Executing Masscan scan", zap.String("target", req.Target))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteAutoRecon executes an AutoRecon scan
func (m *Manager) ExecuteAutoRecon(req models.AutoReconRequest) map[string]interface{} {
	args := []string{req.Target}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("autorecon", args...)
	m.logger.Info("Executing AutoRecon scan", zap.String("target", req.Target))

	result := m.executor.Execute(command, true)
	return m.formatResult(result)
}

// ExecuteMSFVenom executes MSFVenom for payload generation
func (m *Manager) ExecuteMSFVenom(req models.MSFVenomRequest) map[string]interface{} {
	args := []string{"-p", req.Payload}
	if req.Format != "" {
		args = append(args, "-f", req.Format)
	}
	if req.OutputFile != "" {
		args = append(args, "-o", req.OutputFile)
	}
	if req.Encoder != "" {
		args = append(args, "-e", req.Encoder)
	}
	if req.Iterations != "" {
		args = append(args, "-i", req.Iterations)
	}
	if req.AdditionalArgs != "" {
		args = append(args, req.AdditionalArgs)
	}

	command := m.buildCommand("msfvenom", args...)
	m.logger.Info("Executing MSFVenom", zap.String("payload", req.Payload))

	result := m.executor.Execute(command, false) // Don't cache payload generation
	return m.formatResult(result)
}

func (m *Manager) formatResult(result executor.ExecutionResult) map[string]interface{} {
	return map[string]interface{}{
		"success":        result.Success,
		"stdout":         result.Stdout,
		"stderr":         result.Stderr,
		"return_code":    result.ReturnCode,
		"execution_time": result.ExecutionTime,
		"pid":            result.PID,
	}
}
