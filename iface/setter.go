package iface

type IConfigMapSetter interface {
	SetConfigMap(ConfigMap) error
}

type ICapabilityRegistrySetter interface {
	SetCapabilityRegistry(capabilityRegistry ICapabilityRegistry) error
}

type ISetAuthorizerCapabilityMap interface {
	SetAuthorizerCapabilityMap(authorizerMap map[string]IAuthorizer) error
}

type ISetEventTransmitter interface {
	SetEventTransmitter(eventTransmitter IEventTransmitter) error
}
