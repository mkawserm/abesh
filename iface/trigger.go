package iface

import (
	"github.com/mkawserm/abesh/model"
)

type ITrigger interface {
	ICapability
	IStart
	IStop
	ISetEventTransmitter
	IGetEventTransmitter

	AddService(authorizer IAuthorizer,
		authorizationExpression string,
		triggerValues model.ConfigMap,
		service IService) error
}
