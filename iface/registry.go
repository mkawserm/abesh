package iface

type ICapabilityRegistry interface {
	Capability(contractId string) ICapability
}
