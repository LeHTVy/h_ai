package intelligence

// TargetType represents different target types for intelligent analysis
type TargetType string

const (
	TargetTypeWebApplication TargetType = "web_application"
	TargetTypeNetworkHost     TargetType = "network_host"
	TargetTypeAPIEndpoint     TargetType = "api_endpoint"
	TargetTypeCloudService    TargetType = "cloud_service"
	TargetTypeMobileApp       TargetType = "mobile_app"
	TargetTypeBinaryFile      TargetType = "binary_file"
	TargetTypeUnknown         TargetType = "unknown"
)

// TechnologyStack represents common technology stacks
type TechnologyStack string

const (
	TechApache   TechnologyStack = "apache"
	TechNginx    TechnologyStack = "nginx"
	TechIIS      TechnologyStack = "iis"
	TechNodeJS   TechnologyStack = "nodejs"
	TechPHP      TechnologyStack = "php"
	TechPython   TechnologyStack = "python"
	TechJava     TechnologyStack = "java"
	TechDotNet   TechnologyStack = "dotnet"
	TechWordPress TechnologyStack = "wordpress"
	TechDrupal   TechnologyStack = "drupal"
	TechJoomla   TechnologyStack = "joomla"
	TechReact    TechnologyStack = "react"
	TechAngular  TechnologyStack = "angular"
	TechVue      TechnologyStack = "vue"
	TechUnknown  TechnologyStack = "unknown"
)

// TargetProfile represents comprehensive target analysis profile
type TargetProfile struct {
	Target           string                 `json:"target"`
	TargetType       TargetType             `json:"target_type"`
	IPAddresses      []string               `json:"ip_addresses"`
	OpenPorts        []int                  `json:"open_ports"`
	Services         map[int]string         `json:"services"`
	Technologies     []TechnologyStack      `json:"technologies"`
	CMSType          string                 `json:"cms_type,omitempty"`
	CloudProvider    string                 `json:"cloud_provider,omitempty"`
	SecurityHeaders  map[string]string      `json:"security_headers"`
	SSLInfo          map[string]interface{} `json:"ssl_info"`
	Subdomains       []string               `json:"subdomains"`
	Endpoints        []string               `json:"endpoints"`
	AttackSurfaceScore float64              `json:"attack_surface_score"`
	RiskLevel        string                 `json:"risk_level"`
	ConfidenceScore  float64                `json:"confidence_score"`
}

// AttackStep represents individual step in an attack chain
type AttackStep struct {
	Tool                string                 `json:"tool"`
	Parameters          map[string]interface{} `json:"parameters"`
	ExpectedOutcome     string                 `json:"expected_outcome"`
	SuccessProbability  float64                `json:"success_probability"`
	ExecutionTimeEstimate int                  `json:"execution_time_estimate"`
	Dependencies        []string               `json:"dependencies"`
}

// AttackChain represents a sequence of attacks for maximum impact
type AttackChain struct {
	TargetProfile        *TargetProfile `json:"target_profile"`
	Steps                []AttackStep   `json:"steps"`
	SuccessProbability   float64        `json:"success_probability"`
	EstimatedTime        int            `json:"estimated_time"`
	RequiredTools        []string       `json:"required_tools"`
	RiskLevel            string         `json:"risk_level"`
}

// CalculateSuccessProbability calculates overall success probability
func (ac *AttackChain) CalculateSuccessProbability() {
	if len(ac.Steps) == 0 {
		ac.SuccessProbability = 0.0
		return
	}

	prob := 1.0
	for _, step := range ac.Steps {
		prob *= step.SuccessProbability
	}
	ac.SuccessProbability = prob
}

// AddStep adds a step to the attack chain
func (ac *AttackChain) AddStep(step AttackStep) {
	ac.Steps = append(ac.Steps, step)
	
	// Add tool to required tools if not already present
	found := false
	for _, tool := range ac.RequiredTools {
		if tool == step.Tool {
			found = true
			break
		}
	}
	if !found {
		ac.RequiredTools = append(ac.RequiredTools, step.Tool)
	}
	
	ac.EstimatedTime += step.ExecutionTimeEstimate
	ac.CalculateSuccessProbability()
}
