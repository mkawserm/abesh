package echo

import (
	"context"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
)

type Echo struct {
	mValues map[string]string
}

func (e *Echo) Name() string {
	return "abesh_example_echo"
}

func (e *Echo) Version() string {
	return "0.0.1"
}

func (e *Echo) Source() string {
	return "github.com/mkawserm/abesh/example/echo"
}

func (e *Echo) Runtime() string {
	return string(constant.RuntimeNative)
}

func (e *Echo) Category() string {
	return string(constant.CategoryService)
}

func (e *Echo) ContractId() string {
	return "abesh:ex_echo"
}

func (e *Echo) Values() map[string]string {
	return e.mValues
}

func (e *Echo) Setup() error {
	return nil
}

func (e *Echo) SetValues(values map[string]string) error {
	e.mValues = values

	return nil
}

func (e *Echo) New() iface.ICapability {
	return &Echo{}
}

func (e *Echo) Serve(_ context.Context, _ iface.ICapabilityRegistry, _ *model.Event) (*model.Event, error) {
	outputEvent := &model.Event{
		Metadata: &model.Metadata{
			Headers:        map[string]string{"Content-Type": "application/text"},
			ContractIdList: []string{e.ContractId()},
			StatusCode:     200,
			Status:         "OK",
		},

		Data: []byte("echo"),
	}

	return outputEvent, nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&Echo{})
}
