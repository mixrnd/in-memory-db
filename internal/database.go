package internal

import (
	"bufio"
	"bytes"
	"errors"
	"sync"

	"go.uber.org/zap"

	"in-memory-db/internal/compute"
	"in-memory-db/internal/storage"
)

var ErrInternal = errors.New("internal error")

type Database struct {
	storage storage.Engine
	parser  Parser
	logger  *zap.Logger
	wal     Wal
	walMu   sync.Mutex
}

type Wal interface {
	InitRead(f func([]byte) error) error
	Run() error
	Write(query string) error
	Close() error
}

type Parser interface {
	Parse(cmd string) (compute.Query, error)
}

func NewDatabase(storage storage.Engine, parser Parser, logger *zap.Logger, wal Wal) *Database {
	return &Database{
		storage: storage,
		parser:  parser,
		logger:  logger,
		wal:     wal,
	}
}

func (d *Database) Init() {
	if err := d.wal.InitRead(func(data []byte) error {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			query, err := d.parser.Parse(scanner.Text())
			if err != nil {
				return err
			}
			arguments := query.Args()
			switch query.Command() {
			case compute.SetCommand:
				err := d.storage.Set(arguments[0], arguments[1])
				if err != nil {
					return err
				}
			case compute.DelCommand:
				err := d.storage.Del(arguments[0])
				if err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return
	}

	go func() {
		if err := d.wal.Run(); err != nil {
			d.logger.Error("run wall error", zap.Error(err))
			return
		}
	}()
}

func (d *Database) RunQuery(q string) (string, error) {
	d.logger.Debug("handling query", zap.String("query", q))

	query, err := d.parser.Parse(q)
	if err != nil {
		return "", err
	}

	if query.Command() != compute.GetCommand {
		d.walMu.Lock()
		if err = d.wal.Write(query.ToSting()); err != nil {
			d.logger.Error("write to wal", zap.Error(err))
			d.walMu.Unlock()
			return "", err
		}
		d.walMu.Unlock()
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
