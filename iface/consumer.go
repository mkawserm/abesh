package iface

type IConsumer interface {
	ICapability
	IConsumeInputEvent
	IConsumeOutputEvent
}
