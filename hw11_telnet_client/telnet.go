package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var errNilConn = errors.New("connection is nil")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &tcpClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type tcpClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (t *tcpClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}

	t.conn = conn
	return nil
}

func (t *tcpClient) Send() error {
	if t.conn == nil {
		return errNilConn
	}

	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}

	return nil
}

func (t *tcpClient) Receive() error {
	if t.conn == nil {
		return errNilConn
	}

	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("receive: %w", err)
	}

	return nil
}

func (t *tcpClient) Close() error {
	if t.conn == nil {
		return errNilConn
	}

	return t.conn.Close()
}
