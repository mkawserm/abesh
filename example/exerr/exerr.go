package echo

import (
	"context"
	"errors"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
)

var ErrNothing = errors.New("the nothing err")

type ExErr struct {
	mValues map[string]string
}

func (e *ExErr) Name() string {
	return "abesh_example_err"
}

func (e *ExErr) Version() string {
	return "0.0.1"
}

func (e *ExErr) Category() string {
	return string(constant.CategoryService)
}

func (e *ExErr) ContractId() string {
	return "abesh:ex_err"
}

func (e *ExErr) GetConfigMap() model.ConfigMap {
	return e.mValues
}

func (e *ExErr) Setup() error {
	return nil
}

func (e *ExErr) SetConfigMap(values model.ConfigMap) error {
	e.mValues = values

	return nil
}

func (e *ExErr) New() iface.ICapability {
	return &ExErr{}
}

func (e *ExErr) Serve(_ context.Context, input *model.Event) (*model.Event, error) {
	return nil, ErrNothing
}

func init() {
	registry.GlobalRegistry().AddCapability(&ExErr{})
}
