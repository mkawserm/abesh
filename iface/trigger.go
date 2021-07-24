package iface

import "context"

type ITrigger interface {
	ICapability

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	AddService(authorizationHandler AuthorizationHandler,
		authorizationExpression string,
		triggerValues map[string]string,
		capabilityRegistry ICapabilityRegistry,
		service IService) error
}
