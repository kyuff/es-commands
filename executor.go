package commands

import (
	"context"
	"fmt"
	"reflect"

	"github.com/kyuff/es"
)

type State interface {
	es.Handler
}

type Executor[C Command, S State] interface {
	Execute(ctx context.Context, cmd C, state S) ([]es.Content, error)
}

type ExecutorFunc[C Command, S State] func(ctx context.Context, cmd C, state S) ([]es.Content, error)

func (fn ExecutorFunc[C, S]) Execute(ctx context.Context, cmd C, state S) ([]es.Content, error) {
	return fn(ctx, cmd, state)
}

func decorateExecutor[C Command, S State, CMD Command](store Store, entityType string, executor Executor[C, S]) func(ctx context.Context, entityID string, command CMD) error {
	var newStateFunc = newInstance[S]()
	return func(ctx context.Context, entityID string, command CMD) error {
		cmd, ok := Command(command).(C)
		if !ok {
			return fmt.Errorf("command %q is %T, expected %T", command.CommandName(), command, cmd)
		}

		stream := store.Open(ctx, entityType, entityID)
		defer func() {
			_ = stream.Close()
		}()

		var state = newStateFunc()
		err := stream.Project(state)
		if err != nil {
			return err
		}

		events, err := executor.Execute(ctx, cmd, state)
		if err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		err = stream.Write(events...)
		if err != nil {
			return err
		}

		return nil
	}
}

func newInstance[T any]() func() T {
	var (
		t   T
		typ = reflect.TypeOf(t)
	)

	return func() T {
		return reflect.New(typ.Elem()).Interface().(T)
	}
}
