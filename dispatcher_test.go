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

func TestDispatcher(t *testing.T) {
	var (
		newEntityType = func() string {
			return fmt.Sprintf("entiti-type-%d", rand.IntN(100))
		}
		newEntityID = func() string {
			return fmt.Sprintf("entiti-id-%d", rand.IntN(100))
		}
	)

	t.Run("fail with unregistered command", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
		)

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.Error(t, err)
	})

	t.Run("fail with nil command", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			entityType = newEntityType()
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
		)

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return nil, nil
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, nil)

		// assert
		assert.Error(t, err)
	})

	t.Run("fail with the stream project", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			stream     = &StreamMock{}
			entityType = newEntityType()
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
		)

		store.OpenFunc = func(ctx context.Context, entityType string, entityID string) es.Stream {
			return stream
		}
		stream.ProjectFunc = func(handler es.Handler) error {
			return errors.New("project-error")
		}
		stream.CloseFunc = func() error {
			return nil
		}

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return nil, errors.New("executor-error")
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.Error(t, err)
	})

	t.Run("fail with the executor", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			stream     = &StreamMock{}
			entityType = newEntityType()
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
		)

		store.OpenFunc = func(ctx context.Context, entityType string, entityID string) es.Stream {
			return stream
		}
		stream.ProjectFunc = func(handler es.Handler) error {
			return nil
		}
		stream.CloseFunc = func() error {
			return nil
		}

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return nil, errors.New("executor-error")
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.Error(t, err)
	})

	t.Run("fail with the stream write", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			stream     = &StreamMock{}
			entityType = newEntityType()
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
		)

		store.OpenFunc = func(ctx context.Context, entityType string, entityID string) es.Stream {
			return stream
		}
		stream.ProjectFunc = func(handler es.Handler) error {
			return nil
		}
		stream.WriteFunc = func(events ...es.Content) error {
			return errors.New("write-error")
		}
		stream.CloseFunc = func() error {
			return nil
		}

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return []es.Content{&ContentMock{}}, nil
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.Error(t, err)
	})

	t.Run("no write on empty executor result", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			stream     = &StreamMock{}
			entityType = newEntityType()
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
		)

		store.OpenFunc = func(ctx context.Context, entityType string, entityID string) es.Stream {
			return stream
		}
		stream.ProjectFunc = func(handler es.Handler) error {
			return nil
		}
		stream.CloseFunc = func() error {
			return nil
		}

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return nil, nil
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 0, len(stream.WriteCalls()))
	})

	t.Run("successfully call executor", func(t *testing.T) {
		var (
			store      = &StoreMock{}
			stream     = &StreamMock{}
			entityType = newEntityType()
			entityID   = newEntityID()
			dispatcher = commands.NewDispatcher(store)
			events     = []es.Content{&ContentMock{}, &ContentMock{}}
		)

		store.OpenFunc = func(ctx context.Context, entityType string, entityID string) es.Stream {
			return stream
		}
		stream.ProjectFunc = func(handler es.Handler) error {
			return nil
		}
		stream.WriteFunc = func(got ...es.Content) error {
			assert.EqualSlice(t, events, got)
			return nil
		}
		stream.CloseFunc = func() error {
			return nil
		}

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return events, nil
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.NoError(t, err)
		assert.Equal(t, 1, len(stream.WriteCalls()))
	})
}
