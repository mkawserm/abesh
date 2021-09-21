package iface

import "github.com/mkawserm/abesh/model"

type IGetConfigMap interface {
	GetConfigMap() model.ConfigMap
}

type IGetCapabilityRegistry interface {
	GetCapabilityRegistry() ICapabilityRegistry
}

type IGetEventTransmitter interface {
	GetEventTransmitter() IEventTransmitter
}
