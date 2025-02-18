package commands

import (
	"context"
	"fmt"
	"sync"

	"github.com/kyuff/es"
)

type Store interface {
	Open(ctx context.Context, entityType string, entityID string) es.Stream
}

func NewDispatcher(store Store) *Dispatcher {
	return &Dispatcher{
		store:     store,
		executors: make(map[string]func(ctx context.Context, entityID string, cmd Command) error),
	}
}

type Dispatcher struct {
	store     Store
	mux       sync.RWMutex
	executors map[string]func(ctx context.Context, entityID string, cmd Command) error
}

func (d *Dispatcher) Dispatch(ctx context.Context, entityID string, cmd Command) error {
	if cmd == nil {
		return fmt.Errorf("command %T is nil", cmd)
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	executor, ok := d.executors[cmd.Name()]
	if !ok {
		return fmt.Errorf("command %s not registered", cmd.Name())
	}

	return executor(ctx, entityID, cmd)
}
