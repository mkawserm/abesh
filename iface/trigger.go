package iface

type ITrigger interface {
	ICapability
	AddTrigger(triggerValues map[string]string, service IService) error
}
