package iface

import (
	"context"
	"github.com/mkawserm/abesh/model"
)

type IService interface {
	ICapability
	Process(ctx context.Context,
		capabilityRegistry ICapabilityRegistry,
		event *model.Event) *model.Event
}
