package model

type CapabilityManifest struct {
	ContractId    string    `yaml:"contract_id" json:"contract_id"`
	NewContractId string    `yaml:"new_contract_id"`
	Values        ConfigMap `yaml:"values" json:"values"`
}

type TriggerManifest struct {
	Trigger              string    `yaml:"trigger" json:"trigger"`
	TriggerValues        ConfigMap `yaml:"trigger_values" json:"trigger_values"`
	Service              string    `yaml:"service" json:"service"`
	Authorizer           string    `yaml:"authorizer" json:"authorizer"`
	AuthorizerExpression string    `yaml:"authorizer_expression" json:"authorizer_expression"`
}

type RPCManifest struct {
	RPC                  string `yaml:"rpc" json:"rpc"`
	Method               string `yaml:"method" json:"method"`
	Authorizer           string `yaml:"authorizer" json:"authorizer"`
	AuthorizerExpression string `yaml:"authorizer_expression" json:"authorizer_expression"`
}

type ConsumerManifest struct {
	Source string `yaml:"source" json:"source"`
	Sink   string `yaml:"sink" json:"sink"`
}

type Manifest struct {
	Version      string                `yaml:"version" json:"version"` // 1
	Capabilities []*CapabilityManifest `yaml:"capabilities" json:"capabilities"`
	Triggers     []*TriggerManifest    `yaml:"triggers" json:"triggers"`
	RPCS         []*RPCManifest        `yaml:"rpcs" json:"rpcs"`
	Consumers    []*ConsumerManifest   `yaml:"consumers" json:"consumers"`
	Start        []string              `yaml:"start" json:"start"`
}
