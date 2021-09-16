package iface

import "github.com/mkawserm/abesh/model"

type IEventTransmitter interface {
	// TransmitInputEvent transmit input event to the platform data
	// aggregator and event should be respected as read only data
	TransmitInputEvent(contractId string, event *model.Event) error

	// TransmitOutputEvent transmit output event to the platform data
	// aggregator and event should be respected as read only data
	TransmitOutputEvent(contractId string, event *model.Event) error
}
