package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/brutella/dnssd"

	"github.com/FUMCGarland/hvac/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	defer stop()

	// TODO: flag for debugging
	rc, err := load("/etc/relay-module.json")
	if err != nil {
		panic(err.Error())
	}

	if rc.MQTTaddr == "" {
		cto, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		log.Info("discovering controller")

		addFn := func(e dnssd.BrowseEntry) {
			log.Info("found", "ip", e.IPs[0], "port", e.Port, "name", e.Name)
			// why slashes? I don't know
			if !strings.HasPrefix(e.Name, `HVAC\ Controller\ -\ `) {
				return
			}

			// stop after first one
			cancel()
			rc.MQTTaddr = fmt.Sprintf("mqtt://%s:%d", e.IPs[0], e.Port)
		}

		err := dnssd.LookupType(cto, "_mqtt._tcp.local.", addFn, func(dnssd.BrowseEntry) {})
		if err != nil && !(err == context.DeadlineExceeded || err == context.Canceled) {
			panic(err.Error())
		}
	}

	rc.setupGPIO()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		start(ctx, rc)
		stop()
	}()

	<-ctx.Done()
	wg.Wait()
	closeGPIO()
}
