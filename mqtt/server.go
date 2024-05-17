package hvacmqtt

import (
	"fmt"
	"os"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

var inline *mqtt.Server
var root string

// Start initalizes and runs the MQTT server, run it in a go()
func Start(c *hvac.MQTTConfig, done <-chan bool) {
	server := mqtt.New(&mqtt.Options{
		InlineClient: true,      // no need to have a distinct client, inline all our calls
		Logger:       log.Get(), // use the main logger
	})
	root = c.Root
	cmdChan := hvac.GetMQTTChan()

	authData, err := os.ReadFile(c.Auth)
	if err != nil {
		log.Error(err.Error())
		panic(err.Error())
	}

	if err := server.AddHook(new(auth.Hook), &auth.Options{Data: authData}); err != nil {
		log.Error(err.Error())
		panic(err.Error())
	}

	tcp := listeners.NewTCP(listeners.Config{Type: listeners.TypeTCP, ID: c.ID, Address: c.ListenAddr, TLSConfig: nil})
	if err := server.AddListener(tcp); err != nil {
		panic(err.Error())
	}

	go func() {
		if err := server.Serve(); err != nil {
			panic(err.Error())
		}
	}()

	const subscriptionID int = 1

	// subscribe to the topics which relay modules and sensors will update
	sub := fmt.Sprintf("%s/%s/+/%s", root, hvac.PumpsTopic, hvac.CurrentStateEndpoint)
	if err := server.Subscribe(sub, subscriptionID, pumpCallbackFn); err != nil {
		log.Error("subscribe failed", "error", err.Error())
	}

	sub = fmt.Sprintf("%s/%s/+/%s", root, hvac.BlowersTopic, hvac.CurrentStateEndpoint)
	if err := server.Subscribe(sub, subscriptionID, blowerCallbackFn); err != nil {
		log.Error("subscribe failed", "error", err.Error())
	}

	sub = fmt.Sprintf("%s/%s/+/%s", root, hvac.ChillersTopic, hvac.CurrentStateEndpoint)
	if err := server.Subscribe(sub, subscriptionID, chillerCallbackFn); err != nil {
		log.Error("subscribe failed", "error", err.Error())
	}

	sub = fmt.Sprintf("%s/%s/+/%s", root, hvac.RoomsTopic, hvac.TempEndpoint)
	if err := server.Subscribe(sub, subscriptionID, tempCallbackFn); err != nil {
		log.Error("subscribe failed", "error", err.Error())
	}

	sub = fmt.Sprintf("%s/%s/+/%s", root, hvac.RoomsTopic, hvac.HumidityEndpoint)
	if err := server.Subscribe(sub, subscriptionID, humidityCallbackFn); err != nil {
		log.Error("subscribe failed", "error", err.Error())
	}

	inline = server

	for {
		select {
		case cmd := <-cmdChan:
			log.Debug("mqtt got command", "cmd", cmd)
			if err := sendTargetState(cmd.DeviceID, &cmd.Command); err != nil {
				log.Error(err.Error())
			}
		case <-done:
			log.Info("Shutting down MQTT")
			close(cmdChan)
			server.Close()
			return
		}
	}
}
