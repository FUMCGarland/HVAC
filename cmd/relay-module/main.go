package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	defer stop()

	// TODO: flag for debugging
	// TODO: move to /etc/relay-module.json & allow cli flag to specify a different config
	rc, err := load("/home/scot/HVAC/relay-module.json")
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		start(ctx, rc)
		stop()
	}()

	<-ctx.Done()
	wg.Wait()
}
