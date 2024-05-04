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

	tcp := listeners.NewTCP(listeners.Config{listeners.TypeTCP, c.ID, c.ListenAddr, nil})
	if err := server.AddListener(tcp); err != nil {
		panic(err.Error())
	}

	go func() {
		if err := server.Serve(); err != nil {
			panic(err.Error())
		}
	}()

	// subscribe to the topics which relay modules and sensors will update
	sub := fmt.Sprintf("%s/pumps/+/currentstate", root)
	server.Subscribe(sub, 1, pumpCallbackFn)

	sub = fmt.Sprintf("%s/blowers/+/currentstate", root)
	server.Subscribe(sub, 1, blowerCallbackFn)

	sub = fmt.Sprintf("%s/chillers/+/currentstate", root)
	server.Subscribe(sub, 1, blowerCallbackFn)

	sub = fmt.Sprintf("%s/rooms/+/temp", root)
	server.Subscribe(sub, 1, tempCallbackFn)

	inline = server

	for {
		select {
		case cmd := <-cmdChan:
			log.Info("mqtt got command", "cmd", cmd)
			switch cmd.Device.(type) {
			case hvac.PumpID:
				cc := hvac.PumpCommand(cmd.Command)
				SendPumpTargetState(cmd.Device.(hvac.PumpID), &cc)
			case hvac.BlowerID:
				cc := hvac.BlowerCommand(cmd.Command)
				SendBlowerTargetState(cmd.Device.(hvac.BlowerID), &cc)
			case hvac.ChillerID:
				cc := hvac.ChillerCommand(cmd.Command)
				SendChillerTargetState(cmd.Device.(hvac.ChillerID), &cc)
			default:
				log.Error("unknown command device type")
			}
		case <-done:
			log.Info("Shutting down MQTT")
			close(cmdChan)
			server.Close()
			return
		}
	}
}
