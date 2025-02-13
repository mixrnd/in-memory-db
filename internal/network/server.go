package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
	"in-memory-db/internal"
)

const tooManyConnectionsMsg = "too many connections"

type Server struct {
	ctx    context.Context
	db     *internal.Database
	logger *zap.Logger

	address        string
	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	connectionNumber int
	mu               sync.Mutex
}

func NewServer(ctx context.Context, address string, db *internal.Database, logger *zap.Logger, options ...ServerOption) *Server {
	srv := &Server{ctx: ctx, address: address, db: db, logger: logger, connectionNumber: 1}

	for _, o := range options {
		o(srv)
	}

	if srv.connectionNumber <= 0 {
		srv.connectionNumber = 100
	}

	if srv.idleTimeout == 0 {
		srv.idleTimeout = time.Minute * 5
	}

	if srv.bufferSize <= 0 {
		srv.bufferSize = 4096
	}

	return srv
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			s.mu.Lock()
			if s.connectionNumber > s.maxConnections {
				s.mu.Unlock()
				conn.Write([]byte(tooManyConnectionsMsg))
				conn.Close()
				continue
			}
			s.connectionNumber++
			s.mu.Unlock()

			go s.handleUserConnection(conn)
		}
	}()

	<-s.ctx.Done()
	l.Close()

	return nil
}

func (s *Server) handleUserConnection(conn net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("handle panic", zap.Any("panic", v))
		}
		s.mu.Lock()
		s.connectionNumber--
		s.mu.Unlock()

		conn.Close()
	}()

	request := make([]byte, s.bufferSize)
	for {
		if s.idleTimeout != 0 {
			if err := conn.SetDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("set deadline error", zap.Error(err))
				break
			}
		}

		readBytes, err := conn.Read(request)
		if err != nil && !errors.Is(err, io.EOF) {
			s.logger.Error("read user input", zap.Error(err))
			break
		}
		if readBytes == s.bufferSize {
			s.logger.Warn("too big user input", zap.Int("buffer size", s.bufferSize))
			break
		}

		var userResp string
		dbResp, err := s.db.RunQuery(string(request[:readBytes]))
		if err != nil {
			userResp = fmt.Sprintf("error: %s", err.Error())
			s.logger.Error("db error", zap.Error(err))
		} else {
			userResp = dbResp
		}

		if _, err := conn.Write([]byte(userResp)); err != nil {
			s.logger.Warn("write user output", zap.Error(err))
			break
		}
	}
}
