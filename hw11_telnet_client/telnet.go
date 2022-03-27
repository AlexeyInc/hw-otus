package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
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
		"tcp",
		address,
		timeout,
		in,
		out,
		nil,
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
	// scanner := bufio.NewScanner(client.in)

	// if scanner.Scan() {
	// 	msgToServer := scanner.Text()
	// 	// fmt.Println(msgToServer)

	// 	client.conn.Write([]byte(fmt.Sprintf("%s\n", msgToServer)))
	// 	return nil
	// }

	//return errors.New("error: unable to scan input")
}

func (client *Client) Receive() error {
	_, err := io.Copy(client.out, client.conn)

	return err
	// scanner := bufio.NewScanner(client.conn)

	// if scanner.Scan() {
	// 	msgFromSever := scanner.Text()
	// 	fmt.Println(msgFromSever)
	// 	client.out.Write([]byte(fmt.Sprintf("%s\n", msgFromSever)))
	// 	return nil
	// }
	// return errors.New("error: unable to scan msg from server")
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
