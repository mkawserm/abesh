package iface

type IRPC interface {
	ICapability
	IStart
	IStop
	ISetEventTransmitter
	IGetEventTransmitter
	IAddAuthorizer
}
