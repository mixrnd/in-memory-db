package network

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"in-memory-db/internal"
	"in-memory-db/internal/compute"
	inmemory "in-memory-db/internal/storage/in-memory"
)

const (
	testServerAddr = "127.0.0.1:3030"
)

func TestServer_Run(t *testing.T) {
	cancel, server := createServer(5)

	go func() {
		err := server.Run()
		assert.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond) //ждём инициализации сервера

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()

		conn, dialErr := net.Dial("tcp", testServerAddr)
		assert.NoError(t, dialErr)

		conn.Write([]byte("SET config 123"))

		buffer := make([]byte, 4096)
		size, clientErr := conn.Read(buffer)
		require.NoError(t, clientErr)

		assert.Equal(t, "[ok]", string(buffer[:size]))

		conn.Write([]byte("GET config"))

		size, clientErr = conn.Read(buffer)
		require.NoError(t, clientErr)

		assert.Equal(t, "123", string(buffer[:size]))
	}()

	go func() {
		defer wg.Done()

		conn, dialErr := net.Dial("tcp", testServerAddr)
		assert.NoError(t, dialErr)

		conn.Write([]byte("SET key 1"))

		buffer := make([]byte, 4096)
		size, clientErr := conn.Read(buffer)
		require.NoError(t, clientErr)

		assert.Equal(t, "[ok]", string(buffer[:size]))

		conn.Write([]byte("DEL key"))

		size, clientErr = conn.Read(buffer)
		require.NoError(t, clientErr)

		assert.Equal(t, "[ok]", string(buffer[:size]))

		conn.Write([]byte("GET key"))

		size, clientErr = conn.Read(buffer)
		require.NoError(t, clientErr)

		assert.Equal(t, "error: not found", string(buffer[:size]))
	}()

	wg.Wait()

	cancel()
}

func TestServer_RunMaxConn(t *testing.T) {
	cancel, server := createServer(2)

	go func() {
		err := server.Run()
		assert.NoError(t, err)
	}()
	time.Sleep(100 * time.Millisecond) //ждём инициализации сервера

	responses := make([]string, 3)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()

		conn, _ := net.Dial("tcp", testServerAddr)
		conn.Write([]byte("cmd"))
		b := make([]byte, 5012)
		n, _ := conn.Read(b)
		responses[0] = string(b[:n])
	}()

	go func() {
		defer wg.Done()

		conn, _ := net.Dial("tcp", testServerAddr)
		conn.Write([]byte("cmd"))
		b := make([]byte, 5012)
		n, _ := conn.Read(b)

		responses[1] = string(b[:n])
	}()

	go func() {
		defer wg.Done()

		conn, _ := net.Dial("tcp", testServerAddr)
		conn.Write([]byte("cmd"))
		b := make([]byte, 5012)
		n, _ := conn.Read(b)

		responses[2] = string(b[:n])
	}()

	wg.Wait()
	cancel()

	require.Contains(t, responses, "too many connections")
}

func createServer(maxConn int) (context.CancelFunc, *Server) {
	logger := zap.NewNop()
	e := inmemory.NewEngine()
	p := compute.NewParser()
	db := internal.NewDatabase(e, p, logger)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	server := NewServer(ctx, testServerAddr, db, logger,
		WithServerIdleTimeout(time.Minute),
		WithServerBufferSize(4096),
		WithServerMaxConnectionsNumber(maxConn),
	)
	return cancel, server
}
