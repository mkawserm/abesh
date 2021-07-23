package model

type AuthorizationManifest struct {
	Immutable      bool     `yaml:"immutable" json:"immutable"`
	Operator       string   `yaml:"operator" json:"operator"`
	ExpressionList []string `yaml:"expression_list" json:"expression_list"`
}

type TriggerManifest struct {
	Name          string            `yaml:"name" json:"name"`
	Version       string            `yaml:"version" json:"version"`
	ContractId    string            `yaml:"contract_id" json:"contract_id"`
	TriggerValues map[string]string `yaml:"trigger_values" json:"trigger_values"`
}

type CapabilityManifest struct {
	Name       string            `yaml:"name" json:"name"`
	Version    string            `yaml:"version" json:"version"`
	ContractId string            `yaml:"contract_id" json:"contract_id"`
	Source     string            `yaml:"source" json:"source"`
	Runtime    string            `yaml:"runtime" json:"runtime"` // native or wasm
	Category   string            `yaml:"category" json:"category"`
	Values     map[string]string `yaml:"values" json:"values"`
}

type ServiceManifest struct {
	Name          string                 `yaml:"name" json:"name"`
	Version       string                 `yaml:"version" json:"version"`
	ContractId    string                 `yaml:"contract_id" json:"contract_id"`
	Source        string                 `yaml:"source" json:"source"`
	Runtime       string                 `yaml:"runtime" json:"runtime"` // native or wasm
	Authorization *AuthorizationManifest `yaml:"authorization" json:"authorization"`
	Values        map[string]string      `yaml:"values" json:"values"`
	Triggers      []*TriggerManifest     `yaml:"triggers" json:"triggers"`
}

type Manifest struct {
	Version      string                `yaml:"version" json:"version"` // 1
	Capabilities []*CapabilityManifest `yaml:"capabilities" json:"capabilities"`
	Services     []*ServiceManifest    `yaml:"services" json:"services"`
}
