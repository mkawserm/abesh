package iface

type IAddAuthorizer interface {
	AddAuthorizer(method string, contractId string, expression string) error
}
