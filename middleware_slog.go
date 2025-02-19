package commands

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func DefaultSlog() Middleware {
	return SLogMiddleware(slog.Default())
}

func SLogMiddleware(logger *slog.Logger) MiddlewareFunc {
	return func(next func(ctx context.Context, command Command) error) func(ctx context.Context, command Command) error {
		return func(ctx context.Context, command Command) error {
			start := time.Now()
			var (
				err      = next(ctx, command)
				duration = time.Since(start)
			)

			log := logger.WithGroup("commands").With(
				"duration", duration.Milliseconds(),
				"name", command.Name(),
			)
			if err != nil {
				log.ErrorContext(ctx, fmt.Sprintf("[commands] %q executed in %s: %s", command.Name(), duration, err))
				return err
			}

			log.InfoContext(ctx, fmt.Sprintf("[commands] %q executed in %s", command.Name(), duration))

			return nil
		}
	}
}
