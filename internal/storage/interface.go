package storage

import (
	"errors"
)

var ErrNotFound = errors.New("not found")

type Engine interface {
	Set(string, string) error
	Get(string) (string, error)
	Del(string) error
}
