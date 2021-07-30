package iface

import "github.com/mkawserm/abesh/model"

type IEventTransmitter interface {
	TransmitInputEvent(contractId string, event *model.Event) error
	TransmitOutputEvent(contractId string, event *model.Event) error
}
