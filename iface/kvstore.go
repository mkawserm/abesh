package iface

import (
	"context"
	"time"
)

type IKVStore interface {
	ICapability

	// Get the key from store
	Get(ctx context.Context, key string, value interface{}) error

	// Set the key value data to the store with ttl
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	// Delete the key
	Delete(ctx context.Context, key string) error
}
