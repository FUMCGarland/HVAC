package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/dnssd"
	"github.com/FUMCGarland/hvac/log"
	"github.com/FUMCGarland/hvac/mqtt"
	"github.com/FUMCGarland/hvac/rest"
)

func main() {
	done := make(chan bool, 1)
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)

	configPathPtr := flag.String("f", "/etc/hvac.json", "Path to the config file")
	dump := flag.Bool("c", false, "Print the parsed config and exist")
	help := flag.Bool("h", false, "Print the help screen and exit")
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	c, err := hvac.LoadConfig(*configPathPtr)
	if err != nil {
		panic(err.Error())
	}

	if *dump {
		log.Info("config", "c", c)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		rest.Start(c, done)
		log.Info("REST down")
		done <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		hvacmqtt.Start(c.MQTT, done)
		log.Info("MQTT down")
		done <- true
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		hvacdnssd.Start(ctx, c)
		log.Info("DNSSD down")
		cancel()
	}()

	select {
	case <-done: // something called for a shutdown, wait until the rest come down too
		log.Info("Waiting for shutdown")
		cancel()
		wg.Wait()
	case sig := <-sigch:
		log.Info("shutdown requested by signal", "signal", sig)
		done <- true
		cancel()
		wg.Wait()
	}

	// ensure any changes are written out
	c.WriteToStore()
}
