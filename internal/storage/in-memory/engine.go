package inmemory

import (
	"in-memory-db/internal/storage"
)

type Engine struct {
	data map[string]string
}

func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) Set(key string, val string) error {
	e.data[key] = val
	return nil
}

func (e *Engine) Get(key string) (string, error) {
	v, ok := e.data[key]
	if !ok {
		return "", storage.ErrNotFound
	}
	return v, nil
}

func (e *Engine) Del(key string) error {
	delete(e.data, key)
	return nil
}
