package iface

type IConfigMapSetter interface {
	SetConfigMap(ConfigMap) error
}

type ICapabilityRegistrySetter interface {
	SetCapabilityRegistry(capabilityRegistry ICapabilityRegistry) error
}
