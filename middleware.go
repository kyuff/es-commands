package commands

import (
	"context"
)

type Middleware interface {
	Intercept(next func(ctx context.Context, command Command) error) func(ctx context.Context, command Command) error
}
type MiddlewareFunc func(next func(ctx context.Context, command Command) error) func(ctx context.Context, command Command) error

func (fn MiddlewareFunc) Intercept(next func(ctx context.Context, command Command) error) func(ctx context.Context, command Command) error {
	return fn(next)
}

func middlewareExecutor(middlewares []Middleware, inner func(ctx context.Context, entityID string, command Command) error) func(ctx context.Context, entityID string, command Command) error {
	return func(ctx context.Context, entityID string, command Command) error {
		next := func(ctx context.Context, command Command) error {
			return inner(ctx, entityID, command)
		}

		for _, mw := range middlewares {
			next = mw.Intercept(next)
		}

		return next(ctx, command)
	}
}
