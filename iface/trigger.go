package iface

import "context"

type ITrigger interface {
	ICapability

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	AddService(capabilityRegistry ICapabilityRegistry, triggerValues map[string]string, service IService) error
}
