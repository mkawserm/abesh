package authorizer

import (
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
)

type Authorizer struct {
	mValues map[string]string
}

func (e *Authorizer) Name() string {
	return "abesh_example_authorizer"
}

func (e *Authorizer) Version() string {
	return "0.0.1"
}

func (e *Authorizer) Source() string {
	return "github.com/mkawserm/abesh/example/authorizer"
}

func (e *Authorizer) Runtime() string {
	return string(constant.RuntimeNative)
}

func (e *Authorizer) Category() string {
	return string(constant.CategoryAuthorizer)
}

func (e *Authorizer) ContractId() string {
	return "abesh:ex_authorizer"
}

func (e *Authorizer) Values() map[string]string {
	return e.mValues
}

func (e *Authorizer) Setup() error {
	return nil
}

func (e *Authorizer) SetValues(values map[string]string) error {
	e.mValues = values

	return nil
}

func (e *Authorizer) New() iface.ICapability {
	return &Authorizer{}
}

func (e *Authorizer) IsAuthorized(expression string, _ *model.Metadata) bool {
	if expression == "allowAll" {
		return true
	}

	return false
}

func init() {
	registry.GlobalRegistry().AddCapability(&Authorizer{})
}
