package main

import (
	"context"

	"github.com/FUMCGarland/hvac/log"
	"github.com/warthog618/go-gpiocdev"
	"golang.org/x/time/rate"
)

// relay board channel : gpio pin
var pins = []int{5, 6, 13, 16, 19, 20, 21, 26}

const chipname = "gpiochip4"

var gpiochip *gpiocdev.Chip
var gpiorunning bool
var rl *rate.Limiter

func (c *RelayConf) setupGPIO() {
	rl = rate.NewLimiter(2, 0)

	var err error
	gpiochip, err = gpiocdev.NewChip(chipname, gpiocdev.WithConsumer("relay-module"))
	if err != nil {
		log.Error("unable to open GPIO chip, not a raspberry pi?", "err", err.Error())
		return
	}

	for k := range c.Relays {
		log.Debug("initializing relay", "pin", c.Relays[k].Pin)
		l, err := gpiochip.RequestLine(int(c.Relays[k].Pin), gpiocdev.AsOutput(), gpiocdev.AsActiveLow)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		defer l.Close()
		log.Debug("result", "lineinfo", l)
		value, err := l.Value()
		if err != nil {
			log.Error(err.Error())
			continue
		}
		log.Debug("result", "value", value, "pin", c.Relays[k].Pin)
		if err = l.SetValue(0); err != nil {
			log.Error(err.Error())
			continue
		}
	}
	gpiorunning = true
}

func closeGPIO() {
	if !gpiorunning {
		return
	}
	gpiochip.Close()
}

func setRelayState(pin uint8, state bool) error {
	if !gpiorunning {
		log.Info("if gpio were running, pin would be set", "pin", pin, "state", state)
		return nil
	}

	// active low
	value := 1
	if !state {
		value = 0
	}

	rl.Wait(context.Background())

	log.Info("gpio setting relay state", "pin", pin, "value", value)
	l, err := gpiochip.RequestLine(int(pin), gpiocdev.AsOutput(0), gpiocdev.AsActiveLow)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer l.Close()
	log.Debug("result", "lineinfo", l)
	if err := l.SetValue(value); err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
