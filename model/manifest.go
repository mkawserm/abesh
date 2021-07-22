package model

type AuthorizationManifest struct {
	Expression string `yaml:"expression" json:"expression"`
}

type TriggerManifest struct {
	ContractId string `yaml:"contract_id" json:"contract_id"`
	Key        string `yaml:"key" json:"key"`
	Value      string `yaml:"value" json:"value"`
}

type CapabilityManifest struct {
	ContractId string                 `yaml:"contract_id" json:"contract_id"`
	Category   string                 `yaml:"category" json:"category"`
	Values     map[string]interface{} `yaml:"values" json:"values"`
}

type ServiceManifest struct {
	ContractId    string                 `yaml:"contract_id" json:"contract_id"`
	Runtime       string                 `yaml:"runtime" json:"runtime"`
	Authorization *AuthorizationManifest `yaml:"authorization" json:"authorization"`
	Values        map[string]interface{} `yaml:"values" json:"values"`
	Triggers      []*TriggerManifest     `yaml:"triggers" json:"triggers"`
	Capabilities  []*CapabilityManifest  `yaml:"capabilities" json:"capabilities"`
}

type Manifest struct {
	Version  string             `yaml:"version" json:"version"`
	Services []*ServiceManifest `yaml:"services" json:"services"`
}
