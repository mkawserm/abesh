package iface

type IConfigMapGetter interface {
	GetConfigMap() ConfigMap
}

type ICapabilityRegistryGetter interface {
	GetCapabilityRegistry() ICapabilityRegistry
}
