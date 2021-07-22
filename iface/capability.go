package iface

type Capability interface {
	Source() string
	Runtime() string
	Category() string
	ContractId() string
	Values() map[string]interface{}

	SetValues(map[string]interface{})
}
