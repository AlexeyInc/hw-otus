package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Close() error
	Connect() error
	Send() error
	Receive() error
}

type Client struct {
	Network string
	address string
	TimeD   time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		"tcp", address, timeout, in, out, nil,
	}
}

func (client *Client) Connect() error {
	conn, err := net.DialTimeout(client.Network, client.address, client.TimeD)
	if err != nil {
		return err
	}
	client.conn = conn
	return err
}

func (client *Client) Close() error {
	return client.conn.Close()
}

func (client *Client) Send() error {
	_, err := io.Copy(client.conn, client.in)
	return err
}

func (client *Client) Receive() error {
	_, err := io.Copy(client.out, client.conn)
	return err
}
