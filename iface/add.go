package iface

import (
	"embed"
	"github.com/mkawserm/abesh/model"
)

type IAddAuthorizer interface {
	AddAuthorizer(authorizer IAuthorizer, authorizerExpression string, method string) error
}

type IAddService interface {
	AddService(authorizer IAuthorizer,
		authorizerExpression string,
		triggerValues model.ConfigMap,
		service IService) error
}

type IAddEmbeddedStaticFS interface {
	AddEmbeddedStaticFS(pattern string, fs embed.FS)
}
