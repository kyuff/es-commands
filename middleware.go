package commands

import (
	"context"
)

type Middleware[CMD Command] interface {
	Intercept(next func(ctx context.Context, command CMD) error) func(ctx context.Context, command CMD) error
}
type MiddlewareFunc[CMD Command] func(next func(ctx context.Context, command CMD) error) func(ctx context.Context, command CMD) error

func (fn MiddlewareFunc[CMD]) Intercept(next func(ctx context.Context, command CMD) error) func(ctx context.Context, command CMD) error {
	return fn(next)
}

func middlewareExecutor[CMD Command](middlewares []Middleware[CMD], inner func(ctx context.Context, entityID string, command CMD) error) func(ctx context.Context, entityID string, command CMD) error {
	return func(ctx context.Context, entityID string, command CMD) error {
		next := func(ctx context.Context, command CMD) error {
			return inner(ctx, entityID, command)
		}

		for _, mw := range middlewares {
			next = mw.Intercept(next)
		}

		return next(ctx, command)
	}
}
