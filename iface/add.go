package iface

type IAddAuthorizer interface {
	AddAuthorizer(authorizer IAuthorizer, authorizerExpression string, method string) error
}
