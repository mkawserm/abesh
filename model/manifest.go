package model

type AuthorizationManifest struct {
	Expression string `yaml:"expression" json:"expression"`
}

type EnvManifest struct {
	ContractId string                 `yaml:"contract_id" json:"contract_id"`
	Values     map[string]interface{} `yaml:"values" json:"values"`
}

type TriggerManifest struct {
	Trigger string `yaml:"trigger" json:"trigger"`
	Key     string `yaml:"key" json:"key"`
	Value   string `yaml:"value" json:"value"`
}

type CapabilityManifest struct {
	ContractId string `yaml:"contract_id" json:"contract_id"`
}

type ServiceManifest struct {
	ContractId    string                 `yaml:"contract_id" json:"contract_id"`
	Runtime       string                 `yaml:"runtime" json:"runtime"`
	Authorization *AuthorizationManifest `yaml:"authorization" json:"authorization"`

	Triggers     []*TriggerManifest    `yaml:"triggers" json:"triggers"`
	Capabilities []*CapabilityManifest `yaml:"capabilities" json:"capabilities"`
}

type Manifest struct {
	Version  string             `yaml:"version" json:"version"`
	EnvList  []*EnvManifest     `yaml:"env_list" json:"env_list"`
	Services []*ServiceManifest `yaml:"services" json:"services"`
}
