package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

var client *autopaho.ConnectionManager
var stoppedRunning = time.Unix(0, 0)

// start launches the MQTT client, connects to the server, and listens for commands
func start(ctx context.Context, rc *RelayConf) {
	log.Info("Starting up MQTT client ", "on", rc.MQTTaddr)

	u, err := url.Parse(rc.MQTTaddr)
	if err != nil {
		panic(err)
	}

	cliCfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		ConnectUsername:               rc.MQTTuser,
		ConnectPassword:               []byte(rc.MQTTpass),
		KeepAlive:                     20,
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval:         60,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Info("mqtt connection up")
			blowerTopic := fmt.Sprintf("%s/blowers/+/targetstate", rc.Root)
			pumpTopic := fmt.Sprintf("%s/pumps/+/targetstate", rc.Root)
			if _, err := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: blowerTopic, QoS: 1},
					{Topic: pumpTopic, QoS: 1},
				},
			}); err != nil {
				log.Info("failed to subscribe. This is likely to mean no messages will be received.", "err", err.Error())
			}
			log.Info("mqtt subscription made")
		},
		OnConnectError: func(err error) { log.Info("error whilst attempting connection", "err", err.Error()) },
		ClientConfig: paho.ClientConfig{
			ClientID: rc.MQTTclientID,
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				processIncoming,
			},
			OnClientError: func(err error) { log.Error("client error", "err", err.Error()) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Info("server requested disconnect", "reason", d.Properties.ReasonString)
				} else {
					log.Info("server requested disconnect", "reasons code", d.ReasonCode)
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

	log.Info("stopping all relays on restart")
	for k := range rc.Relays {
		if err := setRelayState(rc.Relays[k].Pin, false); err != nil {
			log.Error(err.Error())
			continue
		}
		sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
			CurrentState:  false,
			RanTime:       0,
			TimeRemaining: 0,
		})
	}

	ticker := time.NewTicker(time.Minute)
	curTick := 0
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			// see if any running devices need to stop
			for k := range rc.Relays {
				if rc.Relays[k].Running {
					log.Info("running", "relay", rc.Relays[k].Pin, "StopTime", rc.Relays[k].StopTime, "remaining", rc.Relays[k].StopTime.Sub(now).Minutes())
				}
				if rc.Relays[k].Running && rc.Relays[k].StopTime.Before(now) {
					rt := now.Sub(rc.Relays[k].StartTime)
					log.Info("duration expired", "relay", rc.Relays[k].Pin, "RanTime", rt.Minutes())

					if err := setRelayState(rc.Relays[k].Pin, false); err != nil {
						log.Info(err.Error())
						// don't mark it as stopped so we can retry
					} else {
						rc.Relays[k].Running = false
						rc.Relays[k].RunTime += rt
						rc.Relays[k].StopTime = stoppedRunning
						sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
							CurrentState:  false,
							RanTime:       rt,
							TimeRemaining: 0,
						})
					}
				}
				if curTick%10 == 0 { // send a periodic check-in every 10 minutes
					var tr time.Duration
					if rc.Relays[k].Running { // if rc.Relays[k].StopTime > 0 {
						tr = rc.Relays[k].StopTime.Sub(now)
					}
					sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
						CurrentState:  rc.Relays[k].Running,
						RanTime:       0,
						TimeRemaining: tr,
					})
				}
			}
			curTick++
			continue
		case <-ctx.Done():
		}
		break
	}
	// ctx has already ended
	stopAllRunning(context.Background())

	log.Info("Shutting down MQTT client")
	client.Disconnect(context.Background()) // probably already stopped
}

func stopAllRunning(ctx context.Context) {
	log.Info("Shutting down all relays")
	now := time.Now()
	for k := range rc.Relays {
		if !rc.Relays[k].Running {
			continue
		}
		if err := setRelayState(rc.Relays[k].Pin, false); err != nil {
			log.Error(err.Error())
			continue
		}
		rc.Relays[k].Running = false
		rc.Relays[k].RunTime += now.Sub(rc.Relays[k].StartTime)
		rc.Relays[k].StopTime = stoppedRunning
		sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
			CurrentState:  false,
			RanTime:       now.Sub(rc.Relays[k].StartTime),
			TimeRemaining: 0,
		})
	}
}

func processIncoming(pr paho.PublishReceived) (bool, error) {
	rc := get()

	t := strings.Split(pr.Packet.Topic, "/")
	id, err := strconv.ParseInt(t[len(t)-2], 10, 8)
	if err != nil {
		log.Info("invalid topic object ID", "err", err.Error(), "data", pr.Packet.Topic)
		return false, err
	}
	mode := t[len(t)-3]

	log.Info("processIncoming", "mode", mode, "id", id)

	var cmd hvac.Command
	if err := json.Unmarshal(pr.Packet.Payload, &cmd); err != nil {
		log.Info("unmarshal command failed", "err", err.Error(), "data", pr.Packet.Payload)
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
		log.Debug("request for another controller", "mode", mode, "id", id, "state", cmd.TargetState)
		return true, nil
	}

	log.Info("Toggling Relay", "pin", relay.Pin, "state", cmd.TargetState, "duration", cmd.RunTime.Minutes())
	relay.Running = cmd.TargetState
	if !cmd.TargetState {
		if err := setRelayState(relay.Pin, false); err != nil {
			log.Error(err.Error())
			return true, err
		}
		cmd.RunTime = 0
		relay.StartTime = stoppedRunning
		relay.StopTime = stoppedRunning
	} else {
		switch mode {
		case "pumps":
			if cmd.RunTime < hvac.MinPumpRunTime {
				err := fmt.Errorf("pump runtime too short")
				log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "min", hvac.MinPumpRunTime.Minutes())
				return true, err
			}
			if cmd.RunTime > hvac.MaxPumpRunTime {
				err := fmt.Errorf("pump runtime too long")
				log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "max", hvac.MaxPumpRunTime.Minutes())
				return true, err
			}
		case "blowers":
			if cmd.RunTime < hvac.MinBlowerRunTime {
				err := fmt.Errorf("blower runtime too short")
				log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "min", hvac.MinBlowerRunTime.Minutes())
				return true, err
			}
			if cmd.RunTime > hvac.MaxBlowerRunTime {
				err := fmt.Errorf("blower runtime too long")
				log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "max", hvac.MaxBlowerRunTime.Minutes())
				return true, err
			}
		case "chillers":
			if cmd.RunTime < hvac.MinChillerRunTime {
				err := fmt.Errorf("chiller runtime too short")
				log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "min", hvac.MinChillerRunTime.Minutes())
				return true, err
			}
			if cmd.RunTime > hvac.MaxChillerRunTime {
				err := fmt.Errorf("chiller runtime too long")
				log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "max", hvac.MaxChillerRunTime.Minutes())
				return true, err
			}
		}

		if relay.Running {
			// if currently running, adjust to the stop time that is further out
			newStopTime := time.Now().Add(time.Duration(cmd.RunTime))
			if newStopTime.After(relay.StopTime) {
				log.Info("adjusting stop time", "relay", relay.Pin, "new", newStopTime)
				relay.StopTime = newStopTime
			}
		} else {
			// if not running, record start time and start the relay
			relay.StartTime = time.Now()
			if err := setRelayState(relay.Pin, true); err != nil {
				log.Error(err.Error())
				return true, err
			}
			relay.StopTime = time.Now().Add(time.Duration(cmd.RunTime))
		}
	}

	// Send confirmation
	if err := sendUpdate(context.Background(), rc, relay, &hvac.Response{
		CurrentState:  cmd.TargetState,
		RanTime:       0, // zero on confirmation
		TimeRemaining: relay.StopTime.Sub(time.Now()),
	}); err != nil {
		log.Error(err.Error())
		return false, err
	}
	return true, nil
}

// TODO: add TimeRemaining
func sendUpdate(ctx context.Context, rc *RelayConf, relay *hvac.Relay, response *hvac.Response) error {
	mode := "pumps"
	id := uint8(relay.PumpID)
	if relay.PumpID == 0 {
		mode = "blowers"
		id = uint8(relay.BlowerID)
	}
	topic := fmt.Sprintf("%s/%s/%d/currentstate", rc.Root, mode, id)
	payload, err := json.Marshal(response)
	if err != nil {
		log.Error("json marshal response failed", "err", err.Error())
		return err
	}

	if _, err = client.Publish(ctx, &paho.Publish{
		QoS:     1,
		Topic:   topic,
		Payload: payload,
	}); err != nil {
		if ctx.Err() == nil {
			log.Error("publish response failed", "err", err.Error())
			return err
		}
	}
	return nil
}
