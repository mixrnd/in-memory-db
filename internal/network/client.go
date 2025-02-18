package network

import (
	"errors"
	"io"
	"net"
	"time"
)

var ErrSmallBufferSize = errors.New("small buffer size")

type Client struct {
	address     string
	idleTimeout time.Duration
	bufferSize  int

	conn net.Conn
}

func NewClient(address string, idleTimeout time.Duration, bufferSize int) *Client {
	if idleTimeout == 0 {
		idleTimeout = time.Minute * 5
	}

	if bufferSize <= 0 {
		bufferSize = 4096
	}

	return &Client{address: address, idleTimeout: idleTimeout, bufferSize: bufferSize}
}

func (c *Client) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.address)
	if err != nil {
		return err
	}

	if c.idleTimeout != 0 {
		if err := c.conn.SetDeadline(time.Now().Add(c.idleTimeout)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Send(request string) (string, error) {
	if _, err := c.conn.Write([]byte(request)); err != nil {
		return "", err
	}

	response := make([]byte, c.bufferSize)
	count, err := c.conn.Read(response)
	if err != nil && err != io.EOF {
		return "", err
	} else if count == c.bufferSize {
		return "", ErrSmallBufferSize
	}

	return string(response[:count]), nil
}
