package main

import (
	"github.com/FUMCGarland/hvac/log"
)

func (c *RelayConf) setupGPIO() {
	log.Debug("setupGPIO for macOS")
}

func closeGPIO() {
	log.Debug("closeGPIO for macOS")
}

func setRelayState(pin uint8, state bool) error {
	log.Debug("setRelayState for macOS")
	return nil
}
