package internal

import (
	"errors"

	"go.uber.org/zap"
	"in-memory-db/internal/compute"

	"in-memory-db/internal/storage"
)

var ErrInternal = errors.New("internal error")

type Database struct {
	storage storage.Engine
	parser  Parser
	logger  *zap.Logger
}

type Parser interface {
	Parse(cmd string) (compute.Query, error)
}

func NewDatabase(storage storage.Engine, parser Parser, logger *zap.Logger) *Database {
	return &Database{
		storage: storage,
		parser:  parser,
		logger:  logger,
	}
}

func (d *Database) RunQuery(q string) (string, error) {
	d.logger.Debug("handling query", zap.String("query", q))

	query, err := d.parser.Parse(q)
	if err != nil {
		return "", err
	}

	arguments := query.Args()
	switch query.Command() {
	case compute.GetCommand:
		val, err := d.storage.Get(arguments[0])
		if err != nil {
			return "", err
		}
		return val, nil
	case compute.SetCommand:
		err := d.storage.Set(arguments[0], arguments[1])
		if err != nil {
			return "", err
		}
		return "[ok]", nil
	case compute.DelCommand:
		err := d.storage.Del(arguments[0])
		if err != nil {
			return "", err
		}
		return "[ok]", nil
	}

	d.logger.Error("incorrect query", zap.String("query", q))

	return "internal error", ErrInternal
}
