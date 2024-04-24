package hvacmqtt

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

func blowerCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	ts := strings.Split(pk.TopicName, "/")
	bn, err := strconv.ParseInt(ts[2], 10, 8)
	if err != nil {
		log.Error("invalid blower number", "topic", pk.TopicName, "parsed", bn, "error", err.Error())
		return
	}

	blower := (hvac.BlowerID(bn)).Get()
	if blower == nil {
		log.Error("unknown blower", "blower", bn)
		return
	}

	response := hvac.BlowerResponse{}
	if err := json.Unmarshal(pk.Payload, &response); err != nil {
		log.Error("bad response", "blower", bn, "res", pk.Payload, "err", err.Error())
		return
	}

	log.Info("blower state change", "blower", bn, "state", response.CurrentState)
	blower.Running = response.CurrentState
	if !response.CurrentState && response.RanTime > 0 {
		blower.Runtime += response.RanTime
	}
	if response.CurrentState {
		// now running, log start time
		blower.CurrentStartTime = time.Now()
	} else {
		// now stopped, log stop time
		blower.LastStopTime = time.Now()
		blower.LastStartTime = blower.CurrentStartTime
	}
}

func pumpCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	ts := strings.Split(pk.TopicName, "/")
	pn, err := strconv.ParseInt(ts[2], 10, 8)
	if err != nil {
		log.Error("invalid pump number", "topic", pk.TopicName, "parsed", pn, "error", err.Error())
		return
	}

	pump := (hvac.PumpID(pn)).Get()
	if pump == nil {
		log.Error("unknown pump", "pump", pn)
		return
	}

	response := hvac.PumpResponse{}
	if err := json.Unmarshal(pk.Payload, &response); err != nil {
		log.Error("bad response", "pump", pn, "res", pk.Payload, "err", err.Error())
		return
	}

	log.Info("pump state change", "pump", pn, "state", response.CurrentState)
	pump.Running = response.CurrentState
	if !response.CurrentState && response.RanTime > 0 {
		pump.Runtime += response.RanTime
	}

	if response.CurrentState {
		// now running, log start time
		pump.CurrentStartTime = time.Now()
	} else {
		// now stopped, log stop time
		pump.LastStopTime = time.Now()
		pump.LastStartTime = pump.CurrentStartTime
	}
}

func tempCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	ts := strings.Split(pk.TopicName, "/")
	rn, err := strconv.ParseInt(ts[2], 10, 16)
	if err != nil {
		log.Error("invalid room number", "topic", pk.TopicName, "parsed", rn, "error", err.Error())
		return
	}
	temp, err := strconv.ParseInt(string(pk.Payload), 10, 8)
	if err != nil {
		log.Error("invalid temp", "topic", pk.TopicName, "raw", pk.Payload, "parsed", temp, "error", err.Error())
		return
	}

	room := (hvac.RoomID(rn)).Get()
	if room == nil {
		log.Error("unknown room number", "room", rn)
		return
	}
	log.Info("recording temp", "room", rn, "temp", temp)
	room.Temperature = uint8(temp)
}
