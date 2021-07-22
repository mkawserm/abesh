package iface

type ITrigger interface {
	ICapability
	AddTrigger(key string, value string, service IService) bool
}
