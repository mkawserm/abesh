package iface

type ITrigger interface {
	ICapability

	Start() error
	Stop() error

	AddService(capabilityRegistry ICapabilityRegistry, triggerValues map[string]string, service IService) error
}
