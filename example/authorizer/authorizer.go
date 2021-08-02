package authorizer

import (
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
)

type Authorizer struct {
	mValues iface.ConfigMap
}

func (e *Authorizer) Name() string {
	return "abesh_example_authorizer"
}

func (e *Authorizer) Version() string {
	return "0.0.1"
}

func (e *Authorizer) Category() string {
	return string(constant.CategoryAuthorizer)
}

func (e *Authorizer) ContractId() string {
	return "abesh:ex_authorizer"
}

func (e *Authorizer) GetConfigMap() iface.ConfigMap {
	return e.mValues
}

func (e *Authorizer) Setup() error {
	return nil
}

func (e *Authorizer) SetConfigMap(values iface.ConfigMap) error {
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
