package commands

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/kyuff/es"
)

type Store interface {
	Open(ctx context.Context, entityType string, entityID string) es.Stream
}

func NewDispatcher(store Store, middlewares ...Middleware) *Dispatcher {
	slices.Reverse(middlewares)
	return &Dispatcher{
		store:       store,
		executors:   make(map[string]func(ctx context.Context, entityID string, cmd Command) error),
		middlewares: middlewares,
	}
}

type Dispatcher struct {
	store       Store
	mux         sync.RWMutex
	executors   map[string]func(ctx context.Context, entityID string, cmd Command) error
	middlewares []Middleware
}

func (d *Dispatcher) Dispatch(ctx context.Context, entityID string, cmd Command) error {
	if cmd == nil {
		return fmt.Errorf("command %T is nil", cmd)
	}

	d.mux.RLock()
	defer d.mux.RUnlock()

	executor, ok := d.executors[cmd.CommandName()]
	if !ok {
		return fmt.Errorf("command %s not registered", cmd.CommandName())
	}

	return executor(ctx, entityID, cmd)
}
