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
		newMiddlewareMock = func(n int, calls []int, err error) *MiddlewareMock {
			return &MiddlewareMock{
				InterceptFunc: func(next func(ctx context.Context, command commands.Command) error) func(ctx context.Context, command commands.Command) error {
					return func(ctx context.Context, command commands.Command) error {
						calls[n] = n
						if err != nil {
							return err
						}

						return next(ctx, command)
					}
				},
			}
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

	t.Run("fail with command registered under another name", func(t *testing.T) {
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
		err := dispatcher.Dispatch(t.Context(), entityID, TestDoubleCommand{})

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

	t.Run("execute middleware in order", func(t *testing.T) {
		var (
			store       = &StoreMock{}
			stream      = &StreamMock{}
			entityType  = newEntityType()
			entityID    = newEntityID()
			events      = []es.Content{&ContentMock{}, &ContentMock{}}
			calls       = make([]int, 4)
			middlewares = []*MiddlewareMock{
				newMiddlewareMock(0, calls, nil),
				newMiddlewareMock(1, calls, nil),
				newMiddlewareMock(2, calls, nil),
				newMiddlewareMock(3, calls, nil),
			}
			dispatcher = commands.NewDispatcher(store, middlewares[0], middlewares[1], middlewares[2], middlewares[3])
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
		for _, m := range middlewares {
			assert.Equal(t, 1, len(m.InterceptCalls()))
		}
		assert.EqualSlice(t, []int{0, 1, 2, 3}, calls)
	})

	t.Run("fail with middleware", func(t *testing.T) {
		var (
			store       = &StoreMock{}
			stream      = &StreamMock{}
			entityType  = newEntityType()
			entityID    = newEntityID()
			events      = []es.Content{&ContentMock{}, &ContentMock{}}
			calls       = make([]int, 4)
			middlewares = []*MiddlewareMock{
				newMiddlewareMock(0, calls, nil),
				newMiddlewareMock(1, calls, errors.New("middleware-error")),
				newMiddlewareMock(2, calls, nil),
				newMiddlewareMock(3, calls, nil),
			}
			dispatcher = commands.NewDispatcher(store, middlewares[0], middlewares[1], middlewares[2], middlewares[3])
		)

		_ = commands.RegisterFunc(dispatcher, entityType, func(ctx context.Context, cmd TestCommand, state *StateMock) ([]es.Content, error) {
			return events, nil
		})

		// act
		err := dispatcher.Dispatch(t.Context(), entityID, TestCommand{})

		// assert
		assert.Error(t, err)
		assert.Equal(t, 0, len(stream.WriteCalls()))
		for _, m := range middlewares {
			assert.Equal(t, 1, len(m.InterceptCalls()))
		}
		assert.EqualSlice(t, []int{0, 1, 0, 0}, calls)

	})
}
