package iface

type ICapabilityRegistry interface {
	Capability(contractId string) ICapability
}

type ICapabilityRegistryIterator interface {
	Iterator() map[string]ICapability
}
