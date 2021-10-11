package echo

import (
	"context"
	"errors"
	"fmt"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"go.uber.org/zap"
)

var ErrNothing = errors.New("the nothing err")

type ExPanic struct {
	mValues map[string]string
}

func (e *ExPanic) Name() string {
	return "abesh_example_err"
}

func (e *ExPanic) Version() string {
	return "0.0.1"
}

func (e *ExPanic) Category() string {
	return string(constant.CategoryService)
}

func (e *ExPanic) ContractId() string {
	return "abesh:ex_panic"
}

func (e *ExPanic) GetConfigMap() model.ConfigMap {
	return e.mValues
}

func (e *ExPanic) Setup() error {
	return nil
}

func (e *ExPanic) SetConfigMap(values model.ConfigMap) error {
	e.mValues = values

	return nil
}

func (e *ExPanic) New() iface.ICapability {
	return &ExPanic{}
}

func (e *ExPanic) Serve(_ context.Context, input *model.Event) (event *model.Event, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error from recover")
			panicMsg := fmt.Sprintf("%v", r)
			logger.L(e.ContractId()).Info("recovering from panic")

			// add as much information as possible
			logger.L(e.ContractId()).Error("panic stack trace",
				zap.String("host_name", input.Metadata.GetPath()),
				zap.String("path", input.Metadata.GetPath()),
				zap.String("method", input.Metadata.GetMethod()),
				zap.Any("query", input.Metadata.GetQuery()),
				zap.String("panic_msg", panicMsg))
			return
		}
	}()

	panic("Oh! I am panicking. Hurrah")
	return nil, nil
}

func init() {
	registry.GlobalRegistry().AddCapability(&ExPanic{})
}
