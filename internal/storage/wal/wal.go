package wal

import (
	"context"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Wal struct {
	ctx                  context.Context
	FlushingBatchSize    int
	FlushingBatchTimeout time.Duration
	logger               *zap.Logger

	queryBuffer strings.Builder
	segment     *Segment
	mu          sync.Mutex

	data          chan []byte
	writeWaitChan chan struct{}
}

func NewWal(ctx context.Context, FlushingBatchSize int, FlushingBatchTimeout time.Duration, segment *Segment, logger *zap.Logger) *Wal {
	wal := Wal{
		ctx:                  ctx,
		FlushingBatchSize:    FlushingBatchSize,
		FlushingBatchTimeout: FlushingBatchTimeout,
		queryBuffer:          strings.Builder{},
		segment:              segment,
		data:                 make(chan []byte),
		writeWaitChan:        make(chan struct{}),
		logger:               logger,
	}

	if wal.FlushingBatchSize <= 0 {
		wal.FlushingBatchSize = 100
	}

	if wal.FlushingBatchTimeout <= 0 {
		wal.FlushingBatchTimeout = 200 * time.Millisecond
	}

	return &wal
}

func (w *Wal) Init(f func([]byte) error) error {
	return w.segment.Init(f)
}

func (w *Wal) Run() error {
	ticker := time.NewTicker(w.FlushingBatchTimeout)
	defer ticker.Stop()

	writeCtx, writeCancel := context.WithCancel(context.Background())
	go func() {
		defer func() {
			w.segment.Close()
			close(w.writeWaitChan)
		}()
		for d := range w.data {
			if err := w.segment.Write(d); err != nil {
				w.logger.Error("write segment", zap.Error(err))
			}
			select {
			case <-writeCtx.Done():
				//убедимся что больше нечего записывать
				select {
				case d := <-w.data:
					if err := w.segment.Write(d); err != nil {
						w.logger.Error("write segment", zap.Error(err))
					}
				default:
				}
				return
			default:
			}
		}
	}()

	for {
		select {
		case <-w.ctx.Done():
			w.Flush()
			writeCancel()
			return nil
		case <-ticker.C:
			w.Flush()
		}
	}
}

func (w *Wal) Write(query string) error {
	w.mu.Lock()
	var data []byte
	if w.queryBuffer.Len() > w.FlushingBatchSize {
		data = []byte(w.queryBuffer.String())
		w.queryBuffer.Reset()
	}

	w.queryBuffer.WriteString(query + "\n")
	w.mu.Unlock()

	if len(data) > 0 {
		w.data <- data
	}
	return nil
}

func (w *Wal) Flush() {
	w.mu.Lock()
	data := []byte(w.queryBuffer.String())
	w.queryBuffer.Reset()
	w.mu.Unlock()

	w.data <- data
}

func (w *Wal) Close() error {
	return w.segment.Close()
}

func (w *Wal) WaitWrite() {
	<-w.writeWaitChan
}
