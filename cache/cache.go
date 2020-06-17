package cache

import (
	"context"
)

// Cacher ...
type Cacher interface {
	Set(ctx context.Context, key string, value *Cachable) error
	Get(ctx context.Context, key string, src *Cachable) error
}

// Cachable ...
type Cachable interface {
	SetHash()
	GetHash() string
}
