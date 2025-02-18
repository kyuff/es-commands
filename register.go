package commands

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/kyuff/es"
)

func Register[C Command, S es.Handler](dispatcher *Dispatcher, entityType string, executor Executor[C, S]) (err error) {
	dispatcher.mux.Lock()
	defer dispatcher.mux.Unlock()

	defer func() {
		msg := recover()
		if msg != nil {
			err = errors.Join(err, fmt.Errorf("panic in register: %v", msg))
		}
	}()

	var name = getName[C]()
	if _, ok := dispatcher.executors[name]; ok {
		return fmt.Errorf("command already registered: %s", name)
	}

	dispatcher.executors[name] = decorateExecutor(dispatcher.store, entityType, executor)

	return nil
}

func RegisterFunc[C Command, S es.Handler](dispatcher *Dispatcher, entityType string, executor func(ctx context.Context, cmd C, state S) ([]es.Content, error)) error {
	return Register(dispatcher, entityType, ExecutorFunc[C, S](executor))
}

func getName[C Command]() string {
	typ := reflect.TypeFor[C]()
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
		var c = reflect.New(typ).Interface().(C)
		return c.Name()
	} else {
		var c C
		//goland:noinspection ALL
		return c.Name()
	}
}
