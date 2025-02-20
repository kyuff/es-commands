package commands_test

import (
	"context"
	"testing"

	commands "github.com/kyuff/es-commands"
	"github.com/kyuff/es-commands/internal/assert"
)

func TestValidate(t *testing.T) {
	// arraange
	var (
		called     = false
		nextCalled = false
		sut        = commands.Validate(func(next func(ctx context.Context, t commands.Command) error) func(ctx context.Context, t commands.Command) error {
			return func(ctx context.Context, t commands.Command) error {
				called = true
				return next(ctx, t)
			}
		})
	)
	// act
	err := sut.Intercept(func(ctx context.Context, command commands.Command) error {
		nextCalled = true
		return nil
	})(t.Context(), nil)

	// assert
	assert.NoError(t, err)
	assert.Truef(t, called, "expected command to be called")
	assert.Truef(t, nextCalled, "expected next middleware to be called")

}
