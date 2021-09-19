package iface

type IAddAuthorizer interface {
	AddAuthorizer(method string, expression string, authorizer IAuthorizer) error
}
