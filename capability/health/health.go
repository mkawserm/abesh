package health

import (
	"context"
	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/errors"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/logger"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
	"github.com/mkawserm/abesh/status"
	"github.com/mkawserm/abesh/utility"
	"go.uber.org/zap"
)

var SuccessStatus = status.New(1, "A_H", "OK", map[string]string{})
var FailedStatus = errors.New(1, "A_H_F", "NOT OK", map[string]string{})

type Health struct {
	mCM                 model.ConfigMap
	mCapabilityRegistry iface.ICapabilityRegistry
}

func (h *Health) Name() string {
	return "abesh_health"
}

func (h *Health) Version() string {
	return constant.Version
}

func (h *Health) Category() string {
	return string(constant.CategoryService)
}

func (h *Health) ContractId() string {
	return "abesh:health"
}

func (h *Health) GetConfigMap() model.ConfigMap {
	return h.mCM
}

func (h *Health) SetConfigMap(cm model.ConfigMap) error {
	h.mCM = cm
	return nil
}

func (h *Health) SetCapabilityRegistry(capabilityRegistry iface.ICapabilityRegistry) error {
	h.mCapabilityRegistry = capabilityRegistry
	return nil
}

func (h *Health) Serve(_ context.Context, event *model.Event) (outputEvent *model.Event, outputError error) {
	defer func() {
		r := recover()
		if r != nil {
			logger.L(h.ContractId()).Error("panic information", zap.Any("panic", r))
			outputEvent = utility.JSONErrorEventHTTP(FailedStatus, nil, event.Metadata, h.ContractId())
			outputError = nil
			return
		}
	}()

	outputEvent = utility.JSONSuccessEventHTTP(SuccessStatus, nil, event.Metadata, h.ContractId())
	outputError = nil
	return
}

func (h *Health) New() iface.ICapability {
	return &Health{}
}

func init() {
	registry.GlobalRegistry().AddCapability(&Health{})
}
