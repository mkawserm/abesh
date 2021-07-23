package iface

import "github.com/mkawserm/abesh/model"

type IPlatform interface {
	Setup(manifest *model.Manifest) error
	Run()
}
