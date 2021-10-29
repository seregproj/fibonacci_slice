package memorystorage

import (
	"context"
	"errors"
	"sync"
)

type MemoryCache struct {
	sync.RWMutex
	ar map[uint64]uint64
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		ar: make(map[uint64]uint64),
	}
}

func (m *MemoryCache) Get(ctx context.Context, n uint64) (uint64, error) {
	defer m.RUnlock()

	m.RLock()
	v, ok := m.ar[n]
	if !ok {
		return 0, errors.New("not found")
	}

	return v, nil
}

func (m *MemoryCache) Set(ctx context.Context, k uint64, v uint64) error {
	defer m.Unlock()

	m.Lock()
	if _, ok := m.ar[k]; ok {
		return nil
	}

	m.ar[k] = v

	return nil
}
