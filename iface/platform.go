package iface

import "github.com/mkawserm/abesh/model"

type IPlatform interface {
	Run(manifest *model.Manifest)
}
