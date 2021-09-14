package iface

import (
	"github.com/mkawserm/abesh/model"
)

type IPlatformSetup interface {
	Setup(manifest *model.Manifest) error
}

type IPlatformRun interface {
	Run()
}

type IPlatformTriggerCapabilityGetter interface {
	GetTriggersCapability() map[string]ITrigger
}

type IPlatformAuthorizerCapabilityGetter interface {
	GetAuthorizersCapability() map[string]IAuthorizer
}

type IPlatformConsumersCapabilityGetter interface {
	GetConsumersCapability() map[string]IConsumer
}

type IPlatformCapabilityRegistryGetter interface {
	GetCapabilityRegistry() ICapabilityRegistryIterator
}

type IPlatform interface {
	IPlatformRun
	IPlatformSetup
	IPlatformTriggerCapabilityGetter
	IPlatformAuthorizerCapabilityGetter
	IPlatformConsumersCapabilityGetter
	IPlatformCapabilityRegistryGetter
}
