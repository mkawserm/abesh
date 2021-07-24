package iface

import "github.com/mkawserm/abesh/model"

type IAuthorizer interface {
	ICapability
	IsAuthorized(expression string, metadata *model.Metadata) bool
}
