package inmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"in-memory-db/internal/storage"
)

func TestEngine_SetAndGet(t *testing.T) {
	e := NewEngine()

	val := "val"
	key := "key"

	err := e.Set(key, val)
	assert.NoError(t, err)
	actualVal, err := e.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, actualVal)
}

func TestEngine_GetKeyNotFound(t *testing.T) {
	e := NewEngine()

	val := "val"
	key := "key"

	err := e.Set(key, val)
	assert.NoError(t, err)
	_, err = e.Get("other")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestEngine_Del(t *testing.T) {
	e := NewEngine()

	val := "val"
	key := "key"

	err := e.Set(key, val)
	assert.NoError(t, err)
	actualVal, err := e.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, actualVal)
	err = e.Del(key)
	assert.NoError(t, err)
	_, err = e.Get("other")
	assert.ErrorIs(t, err, storage.ErrNotFound)
}
