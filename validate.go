package commands

import "context"

type Validator[T any] = func(next func(ctx context.Context, t T) error) func(ctx context.Context, t T) error

func Validate(validator Validator[Command]) MiddlewareFunc {
	return validator
}
