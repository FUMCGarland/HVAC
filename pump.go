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
	Runtime          uint64 // seconds
	Hot              bool
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

// CanEnable
// (1) the loop Hot must match the current SystemMode -- no running Hot loops in cooling mode
// (2) at least one blower on the loop must be running in cool mode lest the system freeze over
// (3) don't fast-cycle the pumps, if you stop it, leave it stopped for at least 5(?) minutes
func (p PumpID) CanEnable() error {
	if c.ControlMode == ControlOff {
		err := fmt.Errorf("system off, not starting pump")
		return err
	}

	pump := p.Get()
	if pump.Hot && c.SystemMode != SystemModeHeat {
		err := fmt.Errorf("cannot enable hot pump if system is not heating")
		return err
	}

	// we can enable hot pumps with no blowers for the radiator loops
	// if the loop is a cold loop, check for running blowers
	if !pump.Hot {
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

	if !pump.LastStopTime.Before(time.Now().Add(PumpMinTimeBetweenRuns)) {
		err := fmt.Errorf("pump recently stopped, in hold-down state")
		return err
	}

	return nil
}

func (c *Config) GetPumpFromLoop(id LoopID) PumpID {
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

	if _, err := fd.WriteString(string(j)); err != nil {
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

func (p PumpID) Start(duration uint64, source string) {
	cc := MQTTRequest{
		Device: p,
		Command: Command{
			TargetState: true,
			RunTime:     duration,
			Source:      source,
		},
	}
	cmdChan <- cc
}

func (p PumpID) Stop(source string) {
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
