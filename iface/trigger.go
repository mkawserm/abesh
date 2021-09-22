package iface

type ITrigger interface {
	ICapability
	IStart
	IStop
	ISetEventTransmitter
	IGetEventTransmitter
	IAddService
}
