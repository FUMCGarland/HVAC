package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

type PumpID uint8

type Pump struct {
	CurrentStartTime time.Time
	LastStartTime    time.Time
	LastStopTime     time.Time
	Name             string
	Runtime          time.Duration
	ID               PumpID
	Loop             LoopID
	SystemMode       SystemModeT
	Running          bool
}

func (p PumpID) Get() *Pump {
	for _, k := range c.Pumps {
		if k.ID == p {
			return k
		}
	}
	return nil
}

// canEnable
// (1) the loop SystemMode must match the current SystemMode -- no running Hot loops in cooling mode
// (2) at least one blower on the loop must be running in cool mode lest the system freeze over
// (3) don't fast-cycle the pumps, if you stop it, leave it stopped for at least 5(?) minutes
func (p PumpID) canEnable() error {
	if c.ControlMode == ControlOff {
		err := fmt.Errorf("system off, not starting pump")
		return err
	}

	pump := p.Get()
	if pump.SystemMode != c.SystemMode {
		err := fmt.Errorf("cannot enable pump if in different system mode")
		return err
	}

	// we can enable hot pumps with no blowers for the radiator loops
	// if the loop is a cold loop, check for running blowers
	if pump.SystemMode == SystemModeCool {
		blowerRunning := false
		for _, k := range c.Blowers {
			if pump.Loop == k.ColdLoop && k.Running {
				blowerRunning = true
			}
		}
		if !blowerRunning {
			err := fmt.Errorf("cannot enable cold pump: no blowers on loop are running")
			return err
		}
	}

	if !pump.LastStopTime.Before(time.Now().Add(pumpMinTimeBetweenRuns)) {
		err := fmt.Errorf("pump recently stopped, in hold-down state")
		return err
	}

	maxtemp, hotroom := c.maxTempForPump(pump)
	if c.SystemMode == SystemModeHeat && maxtemp > boilerLockoutTemp {
		log.Info("not starting pump, rooms above boilerLockoutTemp", "boilerLockoutTemp", boilerLockoutTemp, "pump", p, "room", hotroom)
	}

	return nil
}

func (c *Config) getPumpFromLoop(id LoopID) PumpID {
	for _, k := range c.Pumps {
		if k.Loop == id {
			return k.ID
		}
	}
	return PumpID(0)
}

func (c *Config) getHeatZonesFromPump(pump *Pump) []ZoneID {
	var zones []ZoneID

	for _, loop := range c.Loops {
		if pump.Loop == loop.ID && loop.RadiantZone != 0 {
			zones = append(zones, loop.RadiantZone)
		}
	}
	for _, blower := range c.Blowers {
		if blower.HotLoop == pump.Loop {
			zones = append(zones, blower.Zone)
		}
	}
	return zones
}

func (c *Config) maxTempForPump(pump *Pump) (DegF, RoomID) {
	var max DegF
	var hotRoom RoomID

	zones := c.getHeatZonesFromPump(pump)

	for _, z := range zones {
		t, r := z.getMaxTemp()
		if t > max {
			hotRoom = r
			max = t
		}
	}

	return max, hotRoom
}

func (p *Pump) writeToStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("pump-%d.json", p.ID))
	log.Debug("writing pump data", "ID", p.ID, "file", path)

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	j, err := json.Marshal(*p)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if _, err := fd.Write(j); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := fd.Close(); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (p *Pump) readFromStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("pump-%d.json", p.ID))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	var in Pump
	if err = json.Unmarshal(data, &in); err != nil {
		log.Error(err.Error())
		return err
	}

	// config file wins for id/name/hot/loop, only restore values which don't belong in the config
	p.LastStartTime = in.LastStartTime
	p.LastStopTime = in.LastStopTime
	p.Runtime = in.Runtime

	return nil
}

// Start sends the command to the MQTT subsystem to tell the relay-module to start the pump
func (p PumpID) Start(duration time.Duration, source string) error {
	log.Debug("starting pump", "ID", p, "source", source)
	if err := p.canEnable(); err != nil {
		log.Error("cannot enable pump", "id", p, "err", err.Error())
		return err
	}

	if duration < MinPumpRunTime {
		err := fmt.Errorf("duration shorter than minimum: requested %.2f min %.2f", duration.Minutes(), MinPumpRunTime.Minutes())
		return err
	}

	if duration > MaxPumpRunTime {
		err := fmt.Errorf("duration longer than maximum: requested %.2f min %.2f", duration.Minutes(), MaxPumpRunTime.Minutes())
		return err
	}

	cc := MQTTRequest{
		DeviceID: p,
		Command: Command{
			TargetState: true,
			RunTime:     duration,
			Source:      source,
		},
	}
	cmdChan <- cc
	return nil
}

// Stop sends the command to the MQTT subsystem to tell the relay-module to stop the pump
// it stops any chillers on the attached loop if no other pumps are runnning
func (p PumpID) Stop(source string) {
	// if we are the last active blower on the loop, ensure that the pump is shut down
	last := true
	pump := p.Get()
	for _, k := range c.Pumps {
		// if heating no need to stop chiller
		if c.SystemMode == SystemModeHeat {
			last = false
			break
		}

		// skip self
		if k.ID == p {
			continue
		}

		// something else is running, we aren't last
		// TODO the running pump needs to be a cool pump, but that is just paranoia
		if k.Running {
			last = false
			break
		}
	}

	if last {
		if cid := c.GetChillerFromLoop(pump.Loop); cid != 0 {
			if chiller := cid.Get(); chiller.Running {
				chiller.ID.Stop("internal")
			}
		}
	}

	log.Debug("stopping pump", "pumpID", p, "source", source)
	cc := MQTTRequest{
		DeviceID: p,
		Command: Command{
			TargetState: false,
			RunTime:     0,
			Source:      source,
		},
	}
	cmdChan <- cc
}

func (p PumpID) getChiller() ChillerID {
	pump := p.Get()
	if pump == nil {
		log.Warn("requested chiller for invalid pump", "pump ID", p)
		return 0
	}
	return pump.getChiller()
}

func (p *Pump) getChiller() ChillerID {
	for _, k := range c.Chillers {
		for _, l := range k.Loops {
			if l == p.Loop {
				return k.ID
			}
		}
	}
	return ChillerID(0)
}
