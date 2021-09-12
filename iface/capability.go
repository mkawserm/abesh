package iface

type ICapability interface {
	Name() string
	Version() string

	Category() string
	ContractId() string

	New() ICapability
}
