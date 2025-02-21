package wal

import (
	"context"
	"strings"
	"sync"
	"time"
)

type Wal struct {
	ctx                  context.Context
	FlushingBatchSize    int
	FlushingBatchTimeout time.Duration

	queryBuffer strings.Builder
	segment     *Segment
	mu          sync.Mutex
}

func NewWal(ctx context.Context, FlushingBatchSize int, FlushingBatchTimeout time.Duration, segment *Segment) *Wal {
	wal := Wal{
		ctx:                  ctx,
		FlushingBatchSize:    FlushingBatchSize,
		FlushingBatchTimeout: FlushingBatchTimeout,
		queryBuffer:          strings.Builder{},
		segment:              segment,
	}

	if wal.FlushingBatchSize <= 0 {
		wal.FlushingBatchSize = 100
	}

	if wal.FlushingBatchTimeout <= 0 {
		wal.FlushingBatchTimeout = 200 * time.Millisecond
	}

	return &wal
}

func (w *Wal) InitRead(f func([]byte) error) error {
	return w.segment.InitRead(f)
}

func (w *Wal) Run() error {
	defer w.segment.Close()
	ticker := time.NewTicker(w.FlushingBatchTimeout)
	for {
		select {
		case <-w.ctx.Done():
			ticker.Stop()
			return w.Flush()
		case <-ticker.C:
			w.mu.Lock()
			if err := w.segment.Write([]byte(w.queryBuffer.String())); err != nil {
				w.mu.Unlock()
				return err
			}
			w.queryBuffer.Reset()
			w.mu.Unlock()
		}
	}
}

func (w *Wal) Write(query string) error {
	defer w.mu.Unlock()
	w.mu.Lock()
	if w.queryBuffer.Len() > w.FlushingBatchSize {
		if err := w.segment.Write([]byte(w.queryBuffer.String())); err != nil {
			return err
		}
		w.queryBuffer.Reset()
	}

	w.queryBuffer.WriteString(query + "\n")
	return nil
}

func (w *Wal) Flush() error {
	defer w.mu.Unlock()
	w.mu.Lock()

	if err := w.segment.Write([]byte(w.queryBuffer.String())); err != nil {
		return err
	}
	w.queryBuffer.Reset()
	return nil
}

func (w *Wal) Close() error {
	return w.segment.Close()
}
