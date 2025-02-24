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

func NewDispatcher[CMD Command](store Store, middlewares ...Middleware[CMD]) *Dispatcher[CMD] {
	slices.Reverse(middlewares)
	return &Dispatcher[CMD]{
		store:       store,
		executors:   make(map[string]func(ctx context.Context, entityID string, cmd CMD) error),
		middlewares: middlewares,
	}
}

type Dispatcher[CMD Command] struct {
	store       Store
	mux         sync.RWMutex
	executors   map[string]func(ctx context.Context, entityID string, cmd CMD) error
	middlewares []Middleware[CMD]
}

func (d *Dispatcher[CMD]) Dispatch(ctx context.Context, entityID string, cmd CMD) error {
	if any(cmd) == nil {
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
