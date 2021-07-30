package model

type AuthorizerManifest struct {
	ContractId string `yaml:"contract_id" json:"contract_id"`
	Expression string `yaml:"expression" json:"expression"`
}

type TriggerManifest struct {
	ContractId string              `yaml:"contract_id" json:"contract_id"`
	Authorizer *AuthorizerManifest `yaml:"authorizer" json:"authorizer"`
	Values     map[string]string   `yaml:"values" json:"values"`
}

type CapabilityManifest struct {
	ContractId string            `yaml:"contract_id" json:"contract_id"`
	Values     map[string]string `yaml:"values" json:"values"`
}

type ServiceManifest struct {
	ContractId string             `yaml:"contract_id" json:"contract_id"`
	Values     map[string]string  `yaml:"values" json:"values"`
	Triggers   []*TriggerManifest `yaml:"triggers" json:"triggers"`
}

type ConsumerManifest struct {
	Source string `yaml:"source" json:"source"`
	Sink   string `yaml:"sink" json:"sink"`
}

type Manifest struct {
	Version      string                `yaml:"version" json:"version"` // 1
	Capabilities []*CapabilityManifest `yaml:"capabilities" json:"capabilities"`
	Services     []*ServiceManifest    `yaml:"services" json:"services"`
	Consumers    []*ConsumerManifest   `yaml:"consumers" json:"consumers"`
}
