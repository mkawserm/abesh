package model

type AuthorizerManifest struct {
	ContractId string   `yaml:"contract_id" json:"contract_id"`
	Operator   string   `yaml:"operator" json:"operator"`
	Rules      []string `yaml:"rules" json:"rules"`
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
	ContractId string              `yaml:"contract_id" json:"contract_id"`
	Authorizer *AuthorizerManifest `yaml:"authorizer" json:"authorizer"`
	Values     map[string]string   `yaml:"values" json:"values"`
	Triggers   []*TriggerManifest  `yaml:"triggers" json:"triggers"`
}

type Manifest struct {
	Version      string                `yaml:"version" json:"version"` // 1
	Capabilities []*CapabilityManifest `yaml:"capabilities" json:"capabilities"`
	Services     []*ServiceManifest    `yaml:"services" json:"services"`
}
