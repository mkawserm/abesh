package iface

type IConfigMapGetter interface {
	GetConfigMap() ConfigMap
}

type ICapabilityRegistryGetter interface {
	GetCapabilityRegistry() ICapabilityRegistry
}

type IGetEventTransmitter interface {
	GetEventTransmitter() IEventTransmitter
}
