package iface

import "context"

type IStart interface {
	Start(ctx context.Context) error
}
