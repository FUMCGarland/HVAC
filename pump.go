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
	for k := range c.Pumps {
		if c.Pumps[k].ID == p {
			return &c.Pumps[k]
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
		for k := range c.Blowers {
			if pump.Loop == c.Blowers[k].ColdLoop && c.Blowers[k].Running {
				blowerRunning = true
			}
		}
		if !blowerRunning {
			err := fmt.Errorf("cannot enable cold pump no blowers on loop are running")
			return err
		}
	}

	if !pump.LastStopTime.Before(time.Now().Add(pumpMinTimeBetweenRuns)) {
		err := fmt.Errorf("pump recently stopped, in hold-down state")
		return err
	}

	// if locked out, see if we are safet to restart
	if c.BoilerLockout {
		boilerReset := true
		// TODO only check rooms on this pump? ... do we need a per-zone boiler lockout instead of a global?
		// TODO make sure temp reports are recent (1 hour)
		for k := range c.Rooms {
			if c.Rooms[k].Temperature != 0 && c.Rooms[k].Temperature > boilerRecoveryTemp {
				// a room above the reset temp, do not reset
				log.Debug("not unlocking boiler, rooms still above max temp")
				boilerReset = false
			}
		}

		if boilerReset {
			log.Warn("all rooms below recovery temp, unlocking boiler")
			c.BoilerLockout = false
		}
	}

	// still locked out, do not return error so pumps can start
	if c.BoilerLockout {
		err := fmt.Errorf("a room is still too hot, boiler locked out, not starting pump")
		log.Warn(err.Error())
		return nil
	}

	return nil
}

func (c *Config) getPumpFromLoop(id LoopID) PumpID {
	for k := range c.Pumps {
		if c.Pumps[k].Loop == id {
			return c.Pumps[k].ID
		}
	}
	return PumpID(0)
}

func (p *Pump) writeToStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("pump-%d.json", p.ID))

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
	for k := range c.Pumps {
		// if heating no need to stop chiller
		if c.SystemMode == SystemModeHeat {
			last = false
			break
		}

		// skip self
		if c.Pumps[k].ID == p {
			continue
		}

		// something else is running, we aren't last
		// TODO the running pump needs to be a cool pump, but that is just paranoia
		if c.Pumps[k].Running {
			last = false
			break
		}
	}

	if last {
		cid := c.GetChillerFromLoop(pump.Loop)
		chiller := cid.Get()
		if cid != 0 && chiller.Running {
			chiller.ID.Stop("internal")
		}
	}

	log.Debug("stopping pump", "pumpID", p)
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
	for k := range c.Chillers {
		for _, l := range c.Chillers[k].Loops {
			if l == p.Loop {
				return c.Chillers[k].ID
			}
		}
	}
	return ChillerID(0)
}
