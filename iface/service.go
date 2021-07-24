package iface

import (
	"context"
	"github.com/mkawserm/abesh/model"
)

type IService interface {
	ICapability

	// Serve is the entry point of service category capability
	// ctx is immutable
	// capabilityRegistry is immutable
	// event need to respected as immutable
	Serve(ctx context.Context,
		capabilityRegistry ICapabilityRegistry,
		event *model.Event) (*model.Event, error)
}
