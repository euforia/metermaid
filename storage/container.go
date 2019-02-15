package storage

import (
	"errors"
	"sync"

	"github.com/euforia/metermaid/types"
)

var ErrNotFound = errors.New("not found")

type Containers interface {
	Get(string) (types.Container, error)
	Set(types.Container) error
	List() ([]types.Container, error)
}

type InmemContainers struct {
	mu sync.RWMutex
	m  map[string]types.Container
}

func NewInmemContainers() *InmemContainers {
	return &InmemContainers{
		m: make(map[string]types.Container),
	}
}

func (store *InmemContainers) Get(id string) (types.Container, error) {
	store.mu.RLock()
	c, ok := store.m[id]
	store.mu.RUnlock()
	if ok {
		return c, nil
	}
	return c, ErrNotFound
}

func (store *InmemContainers) Set(c types.Container) error {
	store.mu.Lock()
	store.m[c.ID] = c
	store.mu.Unlock()
	return nil
}

func (store *InmemContainers) List() ([]types.Container, error) {
	store.mu.RLock()
	list := make([]types.Container, 0, len(store.m))
	for _, c := range store.m {
		list = append(list, c)
	}
	store.mu.RUnlock()
	return list, nil
}
