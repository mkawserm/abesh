package iface

import "github.com/mkawserm/abesh/model"

type IConsumeInputEvent interface {
	ConsumeInputEvent(contractId string, event *model.Event) error
}

type IConsumeOutputEvent interface {
	ConsumeOutputEvent(contractId string, event *model.Event) error
}
