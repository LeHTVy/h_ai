package models

// CommandRequest represents a generic command execution request
type CommandRequest struct {
	Command  string `json:"command"`
	UseCache bool   `json:"use_cache,omitempty"`
}

// NmapRequest represents an Nmap scan request
type NmapRequest struct {
	Target        string `json:"target"`
	ScanType      string `json:"scan_type,omitempty"`
	Ports         string `json:"ports,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
	UseRecovery   bool   `json:"use_recovery,omitempty"`
}

// NmapAdvancedRequest represents an advanced Nmap scan request
type NmapAdvancedRequest struct {
	Target         string `json:"target"`
	ScanType       string `json:"scan_type,omitempty"`
	Ports          string `json:"ports,omitempty"`
	Timing         string `json:"timing,omitempty"`
	NSEScripts     string `json:"nse_scripts,omitempty"`
	OSDetection    bool   `json:"os_detection,omitempty"`
	VersionDetection bool `json:"version_detection,omitempty"`
	Aggressive     bool   `json:"aggressive,omitempty"`
	Stealth        bool   `json:"stealth,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// MetasploitRequest represents a Metasploit module execution request
type MetasploitRequest struct {
	Module  string            `json:"module"`
	Options map[string]string `json:"options,omitempty"`
}

// GobusterRequest represents a Gobuster scan request
type GobusterRequest struct {
	URL           string `json:"url"`
	Mode          string `json:"mode,omitempty"`
	Wordlist      string `json:"wordlist,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// NucleiRequest represents a Nuclei scan request
type NucleiRequest struct {
	Target        string `json:"target"`
	Templates     string `json:"templates,omitempty"`
	Severity      string `json:"severity,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// SqlmapRequest represents a SQLMap scan request
type SqlmapRequest struct {
	URL           string `json:"url"`
	Data          string `json:"data,omitempty"`
	Cookies       string `json:"cookies,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// HydraRequest represents a Hydra brute force request
type HydraRequest struct {
	Target        string `json:"target"`
	Service       string `json:"service"`
	Username      string `json:"username,omitempty"`
	PasswordList  string `json:"password_list,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// FFufRequest represents an FFuf fuzzing request
type FFufRequest struct {
	URL           string `json:"url"`
	Wordlist      string `json:"wordlist,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// NetexecRequest represents a NetExec request
type NetexecRequest struct {
	Target        string `json:"target"`
	Protocol      string `json:"protocol,omitempty"`
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Hash          string `json:"hash,omitempty"`
	Module        string `json:"module,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// AmassRequest represents an Amass enumeration request
type AmassRequest struct {
	Domain        string `json:"domain"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// MasscanRequest represents a Masscan scan request
type MasscanRequest struct {
	Target        string `json:"target"`
	Ports         string `json:"ports,omitempty"`
	Rate          string `json:"rate,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// AutoReconRequest represents an AutoRecon request
type AutoReconRequest struct {
	Target        string `json:"target"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// MSFVenomRequest represents an MSFVenom payload generation request
type MSFVenomRequest struct {
	Payload       string `json:"payload"`
	Format        string `json:"format,omitempty"`
	OutputFile    string `json:"output_file,omitempty"`
	Encoder       string `json:"encoder,omitempty"`
	Iterations    string `json:"iterations,omitempty"`
	AdditionalArgs string `json:"additional_args,omitempty"`
}

// Intelligence requests
type AnalyzeTargetRequest struct {
	Target       string `json:"target"`
	AnalysisType string `json:"analysis_type,omitempty"`
}

type SelectToolsRequest struct {
	Target    string `json:"target"`
	TargetType string `json:"target_type,omitempty"`
}

type OptimizeParametersRequest struct {
	Tool      string                 `json:"tool"`
	Parameters map[string]interface{} `json:"parameters"`
	Context   map[string]interface{} `json:"context,omitempty"`
}
