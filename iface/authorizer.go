package iface

import "github.com/mkawserm/abesh/model"

type IIsAuthorized interface {
	IsAuthorized(expression string, metadata *model.Metadata) bool
}

type IAuthorizer interface {
	ICapability
	IIsAuthorized
}
