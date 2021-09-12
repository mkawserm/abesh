package iface

type ConfigMap map[string]string

type ICapability interface {
	Name() string
	Version() string

	Category() string
	ContractId() string

	New() ICapability
}
