package commands_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	commands "github.com/kyuff/es-commands"
	"github.com/kyuff/es-commands/internal/assert"
)

func TestSLogMiddleware(t *testing.T) {
	t.Run("log error", func(t *testing.T) {
		// arrange
		var (
			buf        = &bytes.Buffer{}
			logger     = slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
			middleware = commands.SLogMiddleware(logger)
			sut        = middleware.Intercept(func(ctx context.Context, command commands.Command) error {
				return errors.New("test error")
			})
		)

		// act
		err := sut(t.Context(), TestCommand{})

		// assert
		assert.Error(t, err)
		msg := buf.String()
		assert.Match(t, "level=ERROR", msg)
		assert.Match(t, "commands.name=TestCommand", msg)
		assert.Match(t, "commands.duration=", msg)
		t.Log(msg)
	})

	t.Run("log info", func(t *testing.T) {
		// arrange
		var (
			buf        = &bytes.Buffer{}
			logger     = slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{}))
			middleware = commands.SLogMiddleware(logger)
			sut        = middleware.Intercept(func(ctx context.Context, command commands.Command) error {
				return nil
			})
		)

		// act
		err := sut(t.Context(), TestCommand{})

		// assert
		assert.NoError(t, err)
		msg := buf.String()
		assert.Match(t, "level=INFO", msg)
		assert.Match(t, "commands.name=TestCommand", msg)
		assert.Match(t, "commands.duration=", msg)
		t.Log(msg)
	})

	t.Run("log default", func(t *testing.T) {
		// act
		got := commands.DefaultSlog()

		// assert
		if got == nil {
			t.Fatal("got nil")
		}
	})
}
