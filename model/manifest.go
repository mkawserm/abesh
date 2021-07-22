package model

type AuthorizationManifest struct {
	Immutable      bool     `yaml:"immutable" json:"immutable"`
	Operator       string   `yaml:"operator" json:"operator"`
	ExpressionList []string `yaml:"expression_list" json:"expression_list"`
}

type TriggerManifest struct {
	ContractId string `yaml:"contract_id" json:"contract_id"`
	Key        string `yaml:"key" json:"key"`
	Value      string `yaml:"value" json:"value"`
}

type CapabilityManifest struct {
	ContractId string                 `yaml:"contract_id" json:"contract_id"`
	Source     string                 `yaml:"source" json:"source"`
	Runtime    string                 `yaml:"runtime" json:"runtime"` // native or wasm
	Category   string                 `yaml:"category" json:"category"`
	Values     map[string]interface{} `yaml:"values" json:"values"`
}

type ServiceManifest struct {
	ContractId    string                 `yaml:"contract_id" json:"contract_id"`
	Source        string                 `yaml:"source" json:"source"`
	Runtime       string                 `yaml:"runtime" json:"runtime"` // native or wasm
	Authorization *AuthorizationManifest `yaml:"authorization" json:"authorization"`
	Values        map[string]interface{} `yaml:"values" json:"values"`
	Triggers      []*TriggerManifest     `yaml:"triggers" json:"triggers"`
	Capabilities  []*CapabilityManifest  `yaml:"capabilities" json:"capabilities"`
}

type Manifest struct {
	Version  string             `yaml:"version" json:"version"` // 1
	Services []*ServiceManifest `yaml:"services" json:"services"`
}
