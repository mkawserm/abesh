package iface

import "github.com/mkawserm/abesh/model"

type IAddAuthorizer interface {
	AddAuthorizer(authorizer IAuthorizer, authorizerExpression string, method string) error
}

type IAddService interface {
	AddService(authorizer IAuthorizer,
		authorizationExpression string,
		triggerValues model.ConfigMap,
		service IService) error
}
