package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/stianeikeland/go-rpio/v4"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"

	"github.com/FUMCGarland/hvac"
)

var client *autopaho.ConnectionManager

// start launches the MQTT client, connects to the server, and listens for commands
func start(ctx context.Context, rc *RelayConf) {
	rc.Log.Info("Starting up MQTT client ", "on", rc.MQTTaddr)

	u, err := url.Parse(rc.MQTTaddr)
	if err != nil {
		panic(err)
	}

	if err := rpio.Open(); err != nil {
		log.Error(err.Error())
		return
	}
	defer rpio.Close()

	cliCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		ConnectUsername:               rc.MQTTuser,
		ConnectPassword:               []byte(rc.MQTTpass),
		KeepAlive:                     20,
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval:         60,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			rc.Log.Info("mqtt connection up")
			blowerTopic := fmt.Sprintf("%s/blowers/+/targetstate", rc.Root)
			pumpTopic := fmt.Sprintf("%s/pumps/+/targetstate", rc.Root)
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: blowerTopic, QoS: 1},
					{Topic: pumpTopic, QoS: 1},
				},
			}); err != nil {
				rc.Log.Info("failed to subscribe. This is likely to mean no messages will be received.", "err", err.Error())
			}
			rc.Log.Info("mqtt subscription made")
		},
		OnConnectError: func(err error) { rc.Log.Info("error whilst attempting connection", "err", err.Error()) },
		ClientConfig: paho.ClientConfig{
			ClientID: rc.MQTTclientID,
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				processIncoming,
			},
			OnClientError: func(err error) { rc.Log.Error("client error", "err", err.Error()) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					rc.Log.Info("server requested disconnect", "reason", d.Properties.ReasonString)
				} else {
					rc.Log.Info("server requested disconnect", "reasons code", d.ReasonCode)
				}
			},
		},
	}

	client, err = autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		panic(err)
	}

	if err = client.AwaitConnection(ctx); err != nil {
		panic(err)
	}

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now().Unix()
			// see if any running devices need to stop
			for k := range rc.Relays {
				if rc.Relays[k].StopTime > 0 && rc.Relays[k].StopTime < now {
					rc.Log.Info("duration expired", "relay", rc.Relays[k])
					pin := rpio.Pin(rc.Relays[k].Pin)
					pin.Low()
					rc.Relays[k].Running = false
					rc.Relays[k].RunTime += (now - rc.Relays[k].StartTime)
					rc.Relays[k].StopTime = 0
					sendUpdate(rc, &rc.Relays[k], &hvac.Response{
						CurrentState: false,
						RanTime:      uint64(now - rc.Relays[k].StartTime),
					})
				}
			}
			continue
		case <-ctx.Done():
		}
		break
	}

	rc.Log.Info("Shutting down MQTT client")
}

func processIncoming(pr paho.PublishReceived) (bool, error) {
	rc := get()

	t := strings.Split(pr.Packet.Topic, "/")
	id, err := strconv.ParseInt(t[len(t)-2], 10, 8)
	if err != nil {
		rc.Log.Info("invalid topic object ID", "err", err.Error(), "data", pr.Packet.Topic)
		return false, err
	}
	mode := t[len(t)-3]

	rc.Log.Info("about", "mode", mode, "id", id)

	var cmd hvac.Command
	if err := json.Unmarshal(pr.Packet.Payload, &cmd); err != nil {
		rc.Log.Info("unmarshal command failed", "err", err.Error(), "data", pr.Packet.Payload)
		return false, err
	}

	var relay *hvac.Relay
	for k := range rc.Relays {
		if mode == "pumps" && hvac.PumpID(id) == rc.Relays[k].PumpID {
			relay = &rc.Relays[k]
			break
		}
		if mode == "blowers" && hvac.BlowerID(id) == rc.Relays[k].BlowerID {
			relay = &rc.Relays[k]
			break
		}
	}
	if relay == nil {
		// not for us, we are done
		rc.Log.Info("request for another controller", "mode", mode, "id", id, "state", cmd.TargetState)
		return true, nil
	}

	rc.Log.Info("Toggling Relay", "pin", relay.Pin, "state", cmd.TargetState, "duration", cmd.RunTime)
	pin := rpio.Pin(relay.Pin)
	relay.Running = cmd.TargetState
	if !cmd.TargetState {
		cmd.RunTime = 0
		relay.StopTime = 0
		pin.Low()
	} else {
		pin.High()
	}
	relay.StartTime = time.Now().Unix()

	// set the stop time if a duration is set
	if cmd.RunTime > 0 {
		relay.StopTime = time.Now().Add(time.Duration(cmd.RunTime) * time.Second).Unix()
	}

	// Send confirmation
	if err := sendUpdate(rc, relay, &hvac.Response{
		CurrentState: cmd.TargetState,
		RanTime:      0,
	}); err != nil {
		return false, err
	}
	return true, nil
}

func sendUpdate(rc *RelayConf, relay *hvac.Relay, response *hvac.Response) error {
	mode := "pumps"
	id := uint8(relay.PumpID)
	if relay.PumpID == 0 {
		mode = "blowers"
		id = uint8(relay.BlowerID)
	}
	topic := fmt.Sprintf("%s/%s/%d/currentstate", rc.Root, mode, id)
	payload, err := json.Marshal(response)
	if err != nil {
		rc.Log.Error("publish response failed", "err", err.Error())
		return err
	}

	ctx := context.Background()
	if _, err = client.Publish(ctx, &paho.Publish{
		QoS:     1,
		Topic:   topic,
		Payload: payload,
	}); err != nil {
		if ctx.Err() == nil {
			rc.Log.Error("publish response failed", "err", err.Error())
			return err
		}
	}
	return nil
}
