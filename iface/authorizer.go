package iface

import "github.com/mkawserm/abesh/model"

type IAuthorizer interface {
	ICapability
	IsAuthorized(metadata *model.Metadata) bool
}
