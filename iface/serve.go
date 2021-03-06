package iface

import (
	"context"
	"github.com/mkawserm/abesh/model"
)

type IServe interface {
	// Serve is the entry point of service category capability
	// ctx is immutable
	// event need to respected as immutable
	Serve(ctx context.Context, event *model.Event) (*model.Event, error)
}
