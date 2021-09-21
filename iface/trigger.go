package iface

import (
	"context"
	"github.com/mkawserm/abesh/model"
)

type ITrigger interface {
	ICapability

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	GetEventTransmitter() IEventTransmitter
	AddEventTransmitter(eventTransmitter IEventTransmitter) error

	AddService(authorizationHandler AuthorizationHandler,
		authorizationExpression string,
		triggerValues model.ConfigMap,
		service IService) error
}
