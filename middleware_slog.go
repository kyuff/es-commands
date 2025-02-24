package commands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func DefaultSlog[CMD Command]() Middleware[CMD] {
	return SLogMiddleware[CMD](slog.Default())
}

func SLogMiddleware[CMD Command](logger *slog.Logger) MiddlewareFunc[CMD] {
	return func(next func(ctx context.Context, command CMD) error) func(ctx context.Context, command CMD) error {
		return func(ctx context.Context, command CMD) error {
			start := time.Now()
			var (
				err      = next(ctx, command)
				duration = time.Since(start)
			)

			log := logger.WithGroup("commands").With(
				"duration", duration.Milliseconds(),
				"name", command.CommandName(),
			)
			if err != nil {
				log.ErrorContext(ctx, fmt.Sprintf("[commands] %q executed in %s: %s", command.CommandName(), duration, err))
				return err
			}

			log.InfoContext(ctx, fmt.Sprintf("[commands] %q executed in %s", command.CommandName(), duration))

			return nil
		}
	}
}
