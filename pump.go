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
	ID               PumpID
	Name             string
	Loop             LoopID
	Runtime          time.Duration
	SystemMode       SystemModeT
	Running          bool
	CurrentStartTime time.Time
	LastStartTime    time.Time
	LastStopTime     time.Time
}

func (p PumpID) Get() *Pump {
	for k := range c.Pumps {
		if c.Pumps[k].ID == p {
			return &c.Pumps[k]
		}
	}
	return nil
}

// TODO: implement this
var boilerLockout bool

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
		Device: p,
		Command: Command{
			TargetState: true,
			RunTime:     duration,
			Source:      source,
		},
	}
	cmdChan <- cc
	return nil
}

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

	log.Info("stopping pump", "pumpID", p)
	cc := MQTTRequest{
		Device: p,
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
