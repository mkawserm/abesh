package iface

import "context"

type IRPC interface {
	ICapability

	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	GetEventTransmitter() IEventTransmitter
	AddEventTransmitter(eventTransmitter IEventTransmitter) error
}
