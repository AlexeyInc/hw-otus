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

const timeoutFlag = "timeout"

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, timeoutFlag, 2*time.Minute, "Default max time 2min")
	flag.Parse()

	startIndxArgs := 1
	if isFlagPassed(timeoutFlag) {
		startIndxArgs++
	}

	host := os.Args[startIndxArgs]
	port := os.Args[startIndxArgs+1]

	address := net.JoinHostPort(host, port)

	serverDisconnected := make(chan bool)
	defer close(serverDisconnected)

	reader := os.Stdin

	client := NewTelnetClient(address, timeout, reader, &bytes.Buffer{})

	if err := client.Connect(); err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Printf("...Connected to %s:%s\n", host, port)
	defer client.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
	SERVER_KILL:
		for {
			err := client.Receive()
			if err != nil {
				serverDisconnected <- true
				break SERVER_KILL
			}
		}
	}()

	go func() {
	CLIENT_KILL:
		for {
			err := client.Send()
			if err != nil {
				// fmt.Println("...Connection was closed by client")
				break CLIENT_KILL
			}
			select {
			case <-serverDisconnected:
				// fmt.Println("...Connection was closed by peer")
				break CLIENT_KILL
			default:
			}
		}
		wg.Done()
	}()

	wg.Wait()
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
