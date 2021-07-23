package platform

import (
	"fmt"
	"github.com/mkawserm/abesh/model"
)

type One struct {
}

func (o *One) Setup(manifest *model.Manifest) error {
	fmt.Printf("%+v", manifest)

	return nil
}

func (o *One) Run() error {

	return nil
}
