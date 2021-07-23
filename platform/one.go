package platform

import (
	"fmt"
	"github.com/mkawserm/abesh/model"
)

type One struct {
}

func (o *One) Run(manifest *model.Manifest) {
	fmt.Printf("%+v", manifest)
}
