package network

import "time"

type ServerOption func(*Server)

func WithServerIdleTimeout(timeout time.Duration) ServerOption {
	return func(server *Server) {
		server.idleTimeout = timeout
	}
}

func WithServerBufferSize(size int) ServerOption {
	return func(server *Server) {
		server.bufferSize = size
	}
}

func WithServerMaxConnectionsNumber(count int) ServerOption {
	return func(server *Server) {
		server.maxConnections = count
	}
}
