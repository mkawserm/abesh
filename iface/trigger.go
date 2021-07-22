package iface

type ITrigger interface {
	ICapability

	// SetCapabilityRegistry(capabilityRegistry ICapabilityRegistry) error

	Setup() error

	Start() error
	Stop() error

	AddService(triggerValues map[string]string, service IService) error
}
