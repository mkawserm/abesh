package iface

import "context"

type IStop interface {
	Stop(ctx context.Context) error
}
