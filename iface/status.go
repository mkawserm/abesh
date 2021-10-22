package iface

type IStatus interface {
	GetCode() uint32
	GetPrefix() string
	GetMessage() string
	GetParams() map[string]string
}
