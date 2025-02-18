package network

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testForClientServerAddr = "127.0.0.1:3031"
)

func TestClient(t *testing.T) {
	client := NewClient(testForClientServerAddr, 2*time.Second, 12)

	l, err := net.Listen("tcp", testForClientServerAddr)
	assert.NoError(t, err)

	go func() {
		handle := func(conn net.Conn) {
			for {
				buffer := make([]byte, 512)
				size, err := conn.Read(buffer)
				if err != nil {
					break
				}
				conn.Write(buffer[:size])
			}
		}

		for {
			conn, err := l.Accept()
			assert.NoError(t, err)

			go handle(conn)
		}
	}()

	time.Sleep(100 * time.Millisecond) //ждём старат сервера

	err = client.Connect()
	assert.NoError(t, err)

	resp, err := client.Send("hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", resp)

	resp, err = client.Send("1234567891234")
	assert.ErrorIs(t, err, ErrSmallBufferSize)

	client.Close()
}
