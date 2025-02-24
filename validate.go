package commands

import "context"

type Validator[T any] = func(next func(ctx context.Context, t T) error) func(ctx context.Context, t T) error

func Validate[CMD Command](validator Validator[CMD]) MiddlewareFunc[CMD] {
	return validator
}
