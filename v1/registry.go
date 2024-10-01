package cplugin

import (
	"context"
	"fmt"
	"plugin"
	"sync"
)

type registry[T any] struct {
	services map[string]T
	mutex    sync.RWMutex
}

func (r *registry[T]) Wait(ctx context.Context, name string) (t T, ok bool) {
	for {
		select {
		case <-ctx.Done():
			return t, false
		default:
			if s, ok := r.Get(name); ok {
				return s, true
			}
		}
	}
}
func (r *registry[T]) Get(name string) (T, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if s, ok := r.services[name]; ok {
		return s, true
	}
	var t T
	return t, false
}
func (r *registry[T]) Add(name string, command T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.services[name] = command
	return nil
}
func (r *registry[T]) AddByPath(name, path string) (t T, err error) {
	var newPlugin *plugin.Plugin
	newPlugin, err = plugin.Open(path)
	if err != nil {
		return
	}
	var symApp plugin.Symbol
	symApp, err = newPlugin.Lookup("Service")
	if err != nil {
		return
	}
	app, ok := symApp.(T)
	if !ok {
		err = fmt.Errorf("invalid type: %T != %T", t, symApp)
		return
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.services[name] = app
	return app, nil
}

var DefaultRegistry = &registry[any]{}

func NewRegistry[T any]() *registry[T] {
	return &registry[T]{
		mutex:    sync.RWMutex{},
		services: make(map[string]T),
	}
}
