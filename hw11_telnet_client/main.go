package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Invalid number of arguments")
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect: %s\n", err.Error())
		return
	}

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT)

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go readRoutine(client, wg)
	go writeRoutine(client, wg)
	go func() {
		<-ctx.Done()
		os.Exit(1)
	}()

	wg.Wait()
}

func readRoutine(client TelnetClient, wg *sync.WaitGroup) {
	defer wg.Done()

	err := client.Receive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read: %s\n", err.Error())
	}
}

func writeRoutine(client TelnetClient, wg *sync.WaitGroup) {
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close client: %s\n", err.Error())
		}
		wg.Done()
	}()

	err := client.Send()
	if err != nil {
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		return
	}

	fmt.Fprintln(os.Stderr, "...EOF")
}
