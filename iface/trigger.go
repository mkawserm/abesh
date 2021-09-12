package iface

import "context"

type ITrigger interface {
	ICapability

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	GetEventTransmitter() IEventTransmitter
	AddEventTransmitter(eventTransmitter IEventTransmitter) error

	AddService(authorizationHandler AuthorizationHandler,
		authorizationExpression string,
		triggerValues ConfigMap,
		service IService) error
}
