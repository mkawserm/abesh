package model

type AuthorizationManifest struct {
	Immutable      bool     `yaml:"immutable" json:"immutable"`
	Operator       string   `yaml:"operator" json:"operator"`
	ExpressionList []string `yaml:"expression_list" json:"expression_list"`
}

type TriggerManifest struct {
	ContractId    string            `yaml:"contract_id" json:"contract_id"`
	TriggerValues map[string]string `yaml:"trigger_values" json:"trigger_values"`
}

type CapabilityManifest struct {
	ContractId string            `yaml:"contract_id" json:"contract_id"`
	Values     map[string]string `yaml:"values" json:"values"`
}

type ServiceManifest struct {
	ContractId    string                 `yaml:"contract_id" json:"contract_id"`
	Authorization *AuthorizationManifest `yaml:"authorization" json:"authorization"`
	Values        map[string]string      `yaml:"values" json:"values"`
	Triggers      []*TriggerManifest     `yaml:"triggers" json:"triggers"`
}

type Manifest struct {
	Version      string                `yaml:"version" json:"version"` // 1
	Capabilities []*CapabilityManifest `yaml:"capabilities" json:"capabilities"`
	Services     []*ServiceManifest    `yaml:"services" json:"services"`
}
