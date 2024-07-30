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
	log.Debug("blowerCallbackFn", "data", pk.Payload)

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

	response := hvac.Response{}
	if err := json.Unmarshal(pk.Payload, &response); err != nil {
		log.Error("bad response", "blower", bn, "res", pk.Payload, "err", err.Error())
		return
	}

	// ignore the routine check-ins if no change
	if response.CurrentState != blower.Running {
		log.Debug("blower state change", "blower", bn, "state", response.CurrentState)
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

		c := hvac.GetConfig()
		if c.SystemMode == hvac.SystemModeCool {
			// TODO need pump.BlowersRunning()
			blowerRunning := false
			for k := range c.Blowers {
				if blower.ColdLoop == c.Blowers[k].ColdLoop && c.Blowers[k].Running {
					blowerRunning = true
					break
				}
			}
			if !blowerRunning {
				log.Debug("last blower on cool loop stopped, shutting down pump for the loop")
				for k := range c.Pumps {
					if c.Pumps[k].Loop == blower.ColdLoop {
						c.Pumps[k].ID.Stop("internal")
					}
				}
			}
		}
	}
}

func pumpCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	log.Debug("pumpCallbackFn", "data", pk.Payload)

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

	response := hvac.Response{}
	if err := json.Unmarshal(pk.Payload, &response); err != nil {
		log.Error("bad response", "pump", pn, "res", pk.Payload, "err", err.Error())
		return
	}

	// ignore the routine check-ins if no change
	if response.CurrentState != pump.Running {
		log.Debug("pump state change", "pump", pn, "state", response.CurrentState)
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

			// shut down the chiller if no other pumps are running for it
			c := hvac.GetConfig()
			chid := c.GetChillerFromLoop(pump.Loop)
			if chid != 0 && !chid.PumpsRunning() {
				chid.Stop("internal")
			}
		}
	}
}

func chillerCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	log.Debug("chillerCallbackFn", "data", pk.Payload)

	ts := strings.Split(pk.TopicName, "/")
	cn, err := strconv.ParseInt(ts[2], 10, 8)
	if err != nil {
		log.Error("invalid chiller number", "topic", pk.TopicName, "parsed", cn, "error", err.Error())
		return
	}

	chiller := (hvac.ChillerID(cn)).Get()
	if chiller == nil {
		log.Error("unknown chiller", "chiller", cn)
		return
	}

	response := hvac.Response{}
	if err := json.Unmarshal(pk.Payload, &response); err != nil {
		log.Error("bad response", "chiller", cn, "res", pk.Payload, "err", err.Error())
		return
	}

	// ignore the routine check-ins if no change
	if response.CurrentState != chiller.Running {
		log.Debug("chiller state change", "chiller", cn, "state", response.CurrentState)
		chiller.Running = response.CurrentState
		if !response.CurrentState && response.RanTime > 0 {
			chiller.Runtime += response.RanTime
		}

		if response.CurrentState {
			// now running, log start time
			chiller.CurrentStartTime = time.Now()
		} else {
			// now stopped, log stop time
			chiller.LastStopTime = time.Now()
			chiller.LastStartTime = chiller.CurrentStartTime
		}
	}
}

func tempCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	log.Debug("tempCallbackFn", "data", pk.Payload)

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
	log.Debug("recording temp", "room", rn, "temp", temp)
	room.SetTemp(hvac.DegF(temp))
}

func humidityCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	log.Debug("humidityCallbackFn", "data", pk.Payload)

	ts := strings.Split(pk.TopicName, "/")
	rn, err := strconv.ParseInt(ts[2], 10, 16)
	if err != nil {
		log.Error("invalid room number", "topic", pk.TopicName, "parsed", rn, "error", err.Error())
		return
	}
	humidity, err := strconv.ParseInt(string(pk.Payload), 10, 8)
	if err != nil {
		log.Error("invalid humidity", "topic", pk.TopicName, "raw", pk.Payload, "parsed", humidity, "error", err.Error())
		return
	}

	room := (hvac.RoomID(rn)).Get()
	if room == nil {
		log.Error("unknown room number", "room", rn)
		return
	}
	log.Debug("recording humidity", "room", rn, "humidity", humidity)
	room.SetHumidity(uint8(humidity))
}

func shellyCallbackFn(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	ts := strings.Split(pk.TopicName, "/")

	if len(ts) != 4 || ts[1] == "" {
		log.Warn("invalid shelly topic")
	}
	room := hvac.GetRoomIDFromShelly(ts[1])

	switch ts[3] {
	case "temperature":
		if string(pk.Payload) == "null" {
			// periodic update
			return
		}
		temp, err := strconv.ParseFloat(string(pk.Payload), 32)
		if err != nil {
			log.Error(err.Error())
			return
		}
		r := room.Get()
		if r == nil {
			log.Warn("unknown shelly", "shelly", ts[1])
			return
		}
		r.SetTemp(hvac.DegF(temp))
	case "humidity":
		if string(pk.Payload) == "null" {
			// periodic update
			return
		}
		hum, err := strconv.ParseFloat(string(pk.Payload), 32)
		if err != nil {
			log.Error(err.Error())
			return
		}
		r := room.Get()
		if r == nil {
			log.Warn("unknown shelly", "shelly", ts[1])
			return
		}
		r.SetHumidity(uint8(hum))
	case "battery":
		if string(pk.Payload) == "null" {
			// periodic update
			return
		}
		batt, err := strconv.ParseFloat(string(pk.Payload), 32)
		if err != nil {
			log.Error(err.Error())
			return
		}
		r := room.Get()
		if r == nil {
			log.Warn("unknown shelly", "shelly", ts[1])
			return
		}
		r.SetBattery(uint8(batt))
	case "error":
		if string(pk.Payload) != "0" {
			log.Error("error", "shelly", ts[1], "room", room, "data", pk.Payload)
			r := room.Get()
			r.SetBattery(1) // visually show that there is an error
			if r == nil {
				log.Warn("unknown shelly", "shelly", ts[1])
				return
			}
		}
	case "act_reasons":
		if string(pk.Payload) != "[\"sensor\"]" {
			log.Info("act_reason", "shelly", ts[1], "room", room, "data", pk.Payload)
		}
	case "ext_power":
		log.Debug("ext_power", "shelly", ts[1], "room", room, "data", pk.Payload)
		if string(pk.Payload) != "true" {
			return
		}
		r := room.Get()
		if r == nil {
			log.Warn("unknown shelly", "shelly", ts[1])
			return
		}
		r.SetBattery(101)
	}
}
