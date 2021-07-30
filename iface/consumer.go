package iface

import "github.com/mkawserm/abesh/model"

type IConsumer interface {
	ICapability

	ConsumeInputEvent(contractId string, event *model.Event) error
	ConsumeOutputEvent(contractId string, event *model.Event) error
}
