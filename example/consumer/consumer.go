package echo

import (
	"fmt"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
)

type Consumer struct {
	mValues model.ConfigMap
}

func (e *Consumer) Name() string {
	return "abesh_example_event_consumer"
}

func (e *Consumer) Version() string {
	return "0.0.1"
}

func (e *Consumer) Category() string {
	return string(constant.CategoryConsumer)
}

func (e *Consumer) ContractId() string {
	return "abesh:ex_event_consumer"
}

func (e *Consumer) GetConfigMap() model.ConfigMap {
	return e.mValues
}

func (e *Consumer) Setup() error {
	return nil
}

func (e *Consumer) SetConfigMap(values model.ConfigMap) error {
	e.mValues = values

	return nil
}

func (e *Consumer) ConsumeInputEvent(contractId string, event *model.Event) error {
	fmt.Printf("IN: %s - %+v\n", contractId, event)
	return nil
}

func (e *Consumer) ConsumeOutputEvent(contractId string, event *model.Event) error {
	fmt.Printf("OUT: %s - %+v\n", contractId, event)
	return nil
}

func (e *Consumer) New() iface.ICapability {
	return &Consumer{}
}

func init() {
	registry.GlobalRegistry().AddCapability(&Consumer{})
}
