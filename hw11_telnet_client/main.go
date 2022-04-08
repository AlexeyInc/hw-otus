//nolint
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

const timeoutFlag = "timeout"

var defaultTimeout = 10 * time.Second

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, timeoutFlag, defaultTimeout, "Default max time 2min")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("wrong address. Define host and port as args")
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Println(err.Error())
		return
	}
	defer client.Close()

	cntx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		if err := client.Send(); err != nil {
			fmt.Fprintf(os.Stderr, "send error: %v", err)
		} else {
			fmt.Fprintf(os.Stderr, "...EOF")
		}
		cancel()
	}()

	go func() {
		if err := client.Receive(); err != nil {
			fmt.Fprintf(os.Stderr, "receive error: %v", err)
		} else {
			fmt.Fprint(os.Stderr, "...Connection was clodes by perr")
		}
		cancel()
	}()

	<-cntx.Done()
}
