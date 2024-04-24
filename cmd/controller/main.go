package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
	"github.com/FUMCGarland/hvac/mqtt"
	"github.com/FUMCGarland/hvac/rest"
)

func main() {
	done := make(chan bool, 1)
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	// TODO: flag for debugging
	// TODO: move to /etc/hvac.json & allow cli flag to specify a different config
	c, err := hvac.LoadConfig("/home/scot/HVAC/hvac.json")
	if err != nil {
		panic(err.Error())
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		rest.Start(c, done)
		done <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		hvacmqtt.Start(c.MQTT, done)
		done <- true
	}()

	select {
	case <-done: // something called for a shutdown, wait until the rest come down too
		log.Info("Waiting for shutdown")
		wg.Wait()
	case sig := <-sigch:
		log.Info("shutdown requested by signal", "signal", sig)
		done <- true
		wg.Wait()
	}

	// ensure any changes are written out
	c.WriteToStore()
}
