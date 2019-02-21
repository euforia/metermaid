package storage

import (
	"errors"
	"sync"

	"github.com/euforia/metermaid/types"
)

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("not found")
)

// Containers implements a Container storage interface
type Containers interface {
	Get(string) (types.Container, error)
	Set(types.Container) error
	List() ([]types.Container, error)
	Iter(func(types.Container) error) error
}

// InmemContainers implements an in memory Containers interface
type InmemContainers struct {
	mu sync.RWMutex
	m  map[string]types.Container
}

// NewInmemContainers returns a new instance of InmemContainers
func NewInmemContainers() *InmemContainers {
	return &InmemContainers{
		m: make(map[string]types.Container),
	}
}

// Get satisfies the Containers interface
func (store *InmemContainers) Get(id string) (types.Container, error) {
	store.mu.RLock()
	c, ok := store.m[id]
	store.mu.RUnlock()
	if ok {
		return c, nil
	}
	return c, ErrNotFound
}

// Set satisfies the Containers interface
func (store *InmemContainers) Set(c types.Container) error {
	store.mu.Lock()
	store.m[c.ID] = c
	store.mu.Unlock()
	return nil
}

// List satisfies the Containers interface
func (store *InmemContainers) List() ([]types.Container, error) {
	store.mu.RLock()
	list := make([]types.Container, 0, len(store.m))
	for _, c := range store.m {
		list = append(list, c)
	}
	store.mu.RUnlock()
	return list, nil
}

// Iter satisfies the Containers interface
func (store *InmemContainers) Iter(f func(types.Container) error) (err error) {
	store.mu.RLock()
	for _, c := range store.m {
		if err = f(c); err != nil {
			break
		}
	}
	store.mu.RUnlock()
	return err
}
