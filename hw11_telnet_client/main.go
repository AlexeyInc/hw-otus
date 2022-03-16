package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// $ go-telnet --timeout=10s host port
// $ go-telnet mysite.ru 8080
// $ go-telnet --timeout=3s 1.1.1.1 123

func main() {
	host := os.Args[1]
	port := os.Args[2]

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 2*time.Minute, "Default max time 2min")

	address := net.JoinHostPort(host, port)

	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	// defer stop()

	serverDisconnected := make(chan bool)
	//clientDisconnected := make(chan bool)
	defer close(serverDisconnected)
	//defer close(clientDisconnected)

	reader := os.Stdin

	client := NewTelnetClient(address, timeout, reader, &bytes.Buffer{})

	if err := client.Connect(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Printf("...Connected to %s:%s\n", host, port)
	defer client.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
	OUT:
		for {
			err := client.Receive()
			if err != nil {
				serverDisconnected <- true
				break OUT
			}
		}
		wg.Done()
	}()

	go func() {
	OUT:
		for {
			err := client.Send()
			if err != nil {
				fmt.Println("...Connection was closed by client")
				wg.Done()
				break OUT
			}
			select {
			//case <-ctx.Done():
			case <-serverDisconnected:
				fmt.Println("...Connection was closed by peer")
				break OUT
			default:
			}
		}
		wg.Done()
	}()

	wg.Wait()
}
