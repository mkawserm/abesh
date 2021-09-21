package iface

import "github.com/mkawserm/abesh/model"

type ISetConfigMap interface {
	SetConfigMap(model.ConfigMap) error
}

type ISetCapabilityRegistry interface {
	SetCapabilityRegistry(capabilityRegistry ICapabilityRegistry) error
}

type ISetAuthorizerCapabilityMap interface {
	SetAuthorizerCapabilityMap(authorizerMap map[string]IAuthorizer) error
}

type ISetEventTransmitter interface {
	SetEventTransmitter(eventTransmitter IEventTransmitter) error
}
