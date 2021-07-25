package iface

type ICapability interface {
	Name() string
	Version() string

	Category() string
	ContractId() string
	Values() map[string]string

	Setup() error
	SetValues(map[string]string) error

	New() ICapability
}
