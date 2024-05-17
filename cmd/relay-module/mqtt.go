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
			blowerTopic := fmt.Sprintf("%s/%s/+/%s", rc.Root, hvac.BlowersTopic, hvac.TargetStateEndpoint)
			pumpTopic := fmt.Sprintf("%s/%s/+/%s", rc.Root, hvac.PumpsTopic, hvac.TargetStateEndpoint)
			chillerTopic := fmt.Sprintf("%s/%s/+/%s", rc.Root, hvac.ChillersTopic, hvac.TargetStateEndpoint)
			if _, err := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: blowerTopic, QoS: hvac.QoS},
					{Topic: pumpTopic, QoS: hvac.QoS},
					{Topic: chillerTopic, QoS: hvac.QoS},
				},
			}); err != nil {
				log.Info("failed to subscribe. This is likely to mean no messages will be received.", "err", err.Error())
			}
			log.Debug("mqtt subscription made")
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

	// force a stop on all relays, announce this to controller
	log.Debug("stopping all relays on restart")
	for k := range rc.Relays {
		if err := setRelayState(rc.Relays[k].Pin, false); err != nil {
			log.Error(err.Error())
			continue
		}
		if err := sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
			CurrentState:  false,
			RanTime:       0,
			TimeRemaining: 0,
		}); err != nil {
			log.Info("unable to send", "error", err.Error())
		}
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
					log.Debug("running", "relay", rc.Relays[k].Pin, "StopTime", rc.Relays[k].StopTime, "remaining", rc.Relays[k].StopTime.Sub(now).Minutes())
				}
				if rc.Relays[k].Running && rc.Relays[k].StopTime.Before(now) {
					rt := now.Sub(rc.Relays[k].StartTime)
					log.Info("duration expired", "relay", rc.Relays[k].Pin, "RanTime", rt.Minutes())

					if err := setRelayState(rc.Relays[k].Pin, false); err != nil {
						log.Error(err.Error())
						// don't mark it as stopped so we can retry
					} else {
						rc.Relays[k].Running = false
						rc.Relays[k].RunTime += rt
						rc.Relays[k].StopTime = stoppedRunning
						if err := sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
							CurrentState:  false,
							RanTime:       rt,
							TimeRemaining: 0,
						}); err != nil {
							log.Error("unable to send", "error", err.Error())
						}
					}
				}
				if curTick%10 == 0 { // send a periodic check-in every 10 minutes
					var tr time.Duration
					if rc.Relays[k].Running { // if rc.Relays[k].StopTime > 0 {
						tr = rc.Relays[k].StopTime.Sub(now)
					}
					if err := sendUpdate(ctx, rc, &rc.Relays[k], &hvac.Response{
						CurrentState:  rc.Relays[k].Running,
						RanTime:       0,
						TimeRemaining: tr,
					}); err != nil {
						log.Error("unable to send", "error", err.Error())
					}
				}
			}
			curTick++
			continue // skips the break below
		case <-ctx.Done():
			// trigger the break below
		}
		break
	}

	// ctx has already ended, use a new one
	stopAllRunning(context.Background())

	log.Info("Shutting down MQTT client")
	if err := client.Disconnect(context.Background()); err != nil {
		log.Error("disconnect failed", "error", err.Error())
	} // probably already stopped
}

func stopAllRunning(ctx context.Context) {
	log.Debug("shutting down all relays")
	for k := range rc.Relays {
		stopRelay(ctx, &rc.Relays[k])
	}
}

func stopRelay(ctx context.Context, r *hvac.Relay) {
	// Debug
	log.Info("stopping relay", "pin", r.Pin)
	if !r.Running {
		return
	}
	if err := setRelayState(r.Pin, false); err != nil {
		log.Error(err.Error())
		return
	}

	now := time.Now()
	var rantime time.Duration
	if r.StartTime != stoppedRunning {
		rantime = now.Sub(r.StartTime)
	}
	r.StartTime = stoppedRunning
	r.StopTime = stoppedRunning
	r.Running = false
	if err := sendUpdate(ctx, rc, r, &hvac.Response{
		CurrentState:  false,
		RanTime:       rantime,
		TimeRemaining: 0,
	}); err != nil {
		log.Error("unable to send", "error", err.Error())
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

	log.Debug("processIncoming", "topic", pr.Packet.Topic, "data", pr.Packet.Payload)

	var cmd hvac.Command
	if err := json.Unmarshal(pr.Packet.Payload, &cmd); err != nil {
		log.Error("unmarshal command failed", "err", err.Error(), "data", pr.Packet.Payload)
		return false, err
	}

	var relay *hvac.Relay
	for k := range rc.Relays {
		switch {
		case mode == hvac.PumpsTopic && hvac.PumpID(id) == rc.Relays[k].PumpID:
			relay = &rc.Relays[k]
		case mode == hvac.BlowersTopic && hvac.BlowerID(id) == rc.Relays[k].BlowerID:
			relay = &rc.Relays[k]
		case mode == hvac.ChillersTopic && hvac.ChillerID(id) == rc.Relays[k].ChillerID:
			relay = &rc.Relays[k]
		}
	}

	if relay == nil {
		log.Debug("request for another relay-module", "mode", mode, "id", id, "state", cmd.TargetState)
		return true, nil
	}

	log.Debug("Relay-Module: Toggling Relay", "pin", relay.Pin, "state", cmd.TargetState, "duration", cmd.RunTime.Minutes())
	if !cmd.TargetState {
		stopRelay(context.Background(), relay)
		return true, nil
	}

	// all of this is already processed in the controller, do we need to double-check it here?
	// paranoia says yes, if someone is sending commands directly via mqtt instead of through the controller's API
	switch mode {
	case hvac.PumpsTopic:
		if cmd.RunTime < hvac.MinPumpRunTime {
			err := fmt.Errorf("pump runtime too short")
			log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "min", hvac.MinPumpRunTime.Minutes())
			return true, nil
		}
		if cmd.RunTime > hvac.MaxPumpRunTime {
			err := fmt.Errorf("pump runtime too long")
			log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "max", hvac.MaxPumpRunTime.Minutes())
			return true, nil
		}
	case hvac.BlowersTopic:
		if cmd.RunTime < hvac.MinBlowerRunTime {
			err := fmt.Errorf("blower runtime too short")
			log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "min", hvac.MinBlowerRunTime.Minutes())
			return true, nil
		}
		if cmd.RunTime > hvac.MaxBlowerRunTime {
			err := fmt.Errorf("blower runtime too long")
			log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "max", hvac.MaxBlowerRunTime.Minutes())
			return true, nil
		}
	case hvac.ChillersTopic:
		if cmd.RunTime < hvac.MinChillerRunTime {
			err := fmt.Errorf("chiller runtime too short")
			log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "min", hvac.MinChillerRunTime.Minutes())
			return true, nil
		}
		if cmd.RunTime > hvac.MaxChillerRunTime {
			err := fmt.Errorf("chiller runtime too long")
			log.Error(err.Error(), "requested", cmd.RunTime.Minutes(), "max", hvac.MaxChillerRunTime.Minutes())
			return true, nil
		}
	}

	if relay.Running {
		// if currently running, adjust to the stop time that is further out
		// do we need to enforce the max runtimes here? or is the risk of shutting down the chiller whil still running the pumps greater?
		newStopTime := time.Now().Add(time.Duration(cmd.RunTime))
		if newStopTime.After(relay.StopTime) {
			log.Info("adjusting stop time", "relay", relay.Pin, "new", newStopTime)
			relay.StopTime = newStopTime
		}
	} else {
		// if not running, record start time and start the relay
		log.Debug("starting relay")
		relay.StartTime = time.Now()
		if err := setRelayState(relay.Pin, true); err != nil {
			log.Error(err.Error())
			return true, nil
		}
		relay.StopTime = time.Now().Add(time.Duration(cmd.RunTime))

		// now record the running state
		relay.Running = true
	}

	// Send confirmation
	if err := sendUpdate(context.Background(), rc, relay, &hvac.Response{
		CurrentState:  relay.Running,
		RanTime:       0, // zero on confirmation
		TimeRemaining: time.Until(relay.StopTime),
	}); err != nil {
		log.Error(err.Error())
		return false, err
	}
	return true, nil
}

func sendUpdate(ctx context.Context, rc *RelayConf, relay *hvac.Relay, response *hvac.Response) error {
	var mode string
	var id uint8
	switch {
	case relay.BlowerID != 0:
		mode = hvac.BlowersTopic
		id = uint8(relay.BlowerID)
	case relay.ChillerID != 0:
		mode = hvac.ChillersTopic
		id = uint8(relay.ChillerID)
	case relay.PumpID != 0:
		mode = hvac.PumpsTopic
		id = uint8(relay.PumpID)
	default:
		err := fmt.Errorf("unknown device type")
		log.Error(err.Error())
		return err
	}

	topic := fmt.Sprintf("%s/%s/%d/%s", rc.Root, mode, id, hvac.CurrentStateEndpoint)
	payload, err := json.Marshal(response)
	if err != nil {
		log.Error("json marshal response failed", "err", err.Error())
		return err
	}
	log.Debug("sendUpdate", "topic", topic, "payload", payload)

	if _, err = client.Publish(ctx, &paho.Publish{
		QoS:     hvac.QoS,
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
