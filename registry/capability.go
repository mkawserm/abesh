package registry

import "github.com/mkawserm/abesh/iface"

type CapabilityRegistry struct {
	cr map[string]iface.ICapability
}

func (c *CapabilityRegistry) RegisterCapability(capability iface.ICapability) {
	c.cr[capability.ContractId()] = capability
}

func (c *CapabilityRegistry) Capability(contractId string) iface.ICapability {
	if capability, ok := c.cr[contractId]; ok {
		return capability
	}

	return nil
}

func (c *CapabilityRegistry) Iterator() map[string]iface.ICapability {
	return c.cr
}

func NewCapabilityRegistry() *CapabilityRegistry {
	return &CapabilityRegistry{cr: make(map[string]iface.ICapability)}
}
