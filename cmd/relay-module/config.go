package main

import (
	"encoding/json"
	"os"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

type RelayConf struct {
	Root         string
	MQTTuser     string
	MQTTpass     string
	MQTTaddr     string
	MQTTclientID string
	Relays       []hvac.Relay
}

var rc *RelayConf = &RelayConf{
	Root:         "hvac",
	MQTTuser:     "relay1",
	MQTTpass:     "relay1",
	MQTTaddr:     "", // mqtt://127.0.0.1",
	MQTTclientID: "relay-module-1",
	Relays:       []hvac.Relay{},
}

func load(filename string) (*RelayConf, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	// overwrite the defaults with what is in the file
	if err := json.Unmarshal(raw, rc); err != nil {
		panic(err.Error())
	}

	log.Start()

	if err := relayvalidate(); err != nil {
		log.Debug("config", "config", rc)
		panic(err.Error())
	}

	return rc, nil
}

func relayvalidate() error {
	// TODO do something
	return nil
}
