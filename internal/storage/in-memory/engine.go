package inmemory

import (
	"sync"

	"in-memory-db/internal/storage"
)

type Engine struct {
	data map[string]string
	mu   sync.Mutex
}

func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

func (e *Engine) Set(key string, val string) error {
	defer e.mu.Unlock()
	e.mu.Lock()
	e.data[key] = val
	return nil
}

func (e *Engine) Get(key string) (string, error) {
	defer e.mu.Unlock()
	e.mu.Lock()
	v, ok := e.data[key]
	if !ok {
		return "", storage.ErrNotFound
	}
	return v, nil
}

func (e *Engine) Del(key string) error {
	defer e.mu.Unlock()
	e.mu.Lock()
	delete(e.data, key)
	return nil
}
