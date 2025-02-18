package commands_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/kyuff/es"
	commands "github.com/kyuff/es-commands"
	"github.com/kyuff/es-commands/internal/assert"
)

func TestRegister(t *testing.T) {
	var (
		newEntityType = func() string {
			return fmt.Sprintf("entiti-type-%d", rand.IntN(100))
		}
	)

	var registerTestCases = []struct {
		name      string
		act       func(t *testing.T, store commands.Store, dispatcher *commands.Dispatcher, entityType string) error
		expectErr bool
	}{

		{
			name: "value command and receiver",
			act: func(t *testing.T, store commands.Store, dispatcher *commands.Dispatcher, entityType string) error {
				return commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
					return nil, nil
				})
			},
			expectErr: false,
		},
		{
			name: "pointer command value receiver",
			act: func(t *testing.T, store commands.Store, dispatcher *commands.Dispatcher, entityType string) error {
				return commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd *TestCommand, state *StateMock) ([]es.Content, error) {
					return nil, nil
				})
			},
			expectErr: false,
		},
		{
			name: "pointer command and receiver",
			act: func(t *testing.T, store commands.Store, dispatcher *commands.Dispatcher, entityType string) error {
				return commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd *TestPointerCommand, state *StateMock) ([]es.Content, error) {
					return nil, nil
				})
			},
			expectErr: false,
		},
		{
			name: "panic command name",
			act: func(t *testing.T, store commands.Store, dispatcher *commands.Dispatcher, entityType string) error {
				return commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestPanicCommand, state *StateMock) ([]es.Content, error) {
					return nil, nil
				})
			},
			expectErr: true,
		},
		{
			name: "register twice",
			act: func(t *testing.T, store commands.Store, dispatcher *commands.Dispatcher, entityType string) error {
				return errors.Join(
					commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
						return nil, nil
					}),
					commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
						return nil, nil
					}),
				)
			},
			expectErr: true,
		},
	}

	for _, tt := range registerTestCases {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			var (
				store      = &StoreMock{}
				dispatcher = commands.NewDispatcher(store)
				entityType = newEntityType()
			)

			// act
			err := tt.act(t, store, dispatcher, entityType)

			// assert
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
