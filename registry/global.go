package registry

import (
	"github.com/mkawserm/abesh/iface"
	"sync"
)

var (
	globalRegistryOnce sync.Once
	globalRegistryIns  *globalRegistry
)

type globalRegistry struct {
	mCapability map[string]iface.ICapability
}

func (g *globalRegistry) setup() {
	g.mCapability = make(map[string]iface.ICapability)
}

func (g *globalRegistry) AddCapability(capability iface.ICapability) {
	g.mCapability[capability.ContractId()] = capability
}

func (g *globalRegistry) GetCapability(contractId string) iface.ICapability {
	return g.mCapability[contractId]
}

func (g *globalRegistry) CapabilityIterator() map[string]iface.ICapability {
	return g.mCapability
}

func GlobalRegistry() *globalRegistry {
	return globalRegistryIns
}

func init() {
	globalRegistryOnce.Do(func() {
		globalRegistryIns = &globalRegistry{}
		globalRegistryIns.setup()
	})
}
