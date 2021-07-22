package iface

type ICapability interface {
	Source() string
	Runtime() string
	Category() string
	ContractId() string
	Values() map[string]string

	SetValues(map[string]string) error
}
