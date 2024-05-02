package main

import (
	"github.com/FUMCGarland/hvac/log"

	"github.com/warthog618/go-gpiocdev"
)

// relay board channel : gpio pin
var pins = []uint8{5, 6, 13, 16, 19, 20, 21, 26}

const chipname = "gpiochip0"

var gpiochip *gpiocdev.Chip
var gpiorunning bool

func setupGPIO() {
	var err error
	gpiochip, err = gpiocdev.NewChip(chipname, gpiocdev.WithConsumer("relay-module"))
	if err != nil {
		log.Error("unable to open GPIO chip, not a raspberry pi?", "err", err.Error())
		//panic(err.Error())
		return
	}

	// initialize -- this just displays info and tests our ability to read from the pins
	for k, v := range pins {
		log.Info("initializing relay", "relay", k, "pin", v)
		l, err := gpiochip.RequestLine(int(v), gpiocdev.AsOutput(0))
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
		log.Info("result", "value", value, "relay", k, "pin", v)
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

	value := 0
	if state {
		value = 1
	}

	log.Info("setting relay state", "pin", pin)
	l, err := gpiochip.RequestLine(int(pin), gpiocdev.AsOutput(0))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer l.Close()
	log.Debug("result", "lineinfo", l)
	err = l.SetValue(value)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
