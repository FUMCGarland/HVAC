package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

// BlowerID a the unique identifier of each individual blower
type BlowerID uint8

// A blower is a device that blows air over a coil to cool/heat the air
// blowers are connected to loops, which are driven by pumps
// a zone is the region of the building serviced by one or more blowers
type Blower struct {
	CurrentStartTime time.Time
	LastStartTime    time.Time
	LastStopTime     time.Time
	Name             string
	Runtime          time.Duration
	FilterTime       uint64
	ID               BlowerID
	HotLoop          LoopID
	ColdLoop         LoopID
	Zone             ZoneID
	Running          bool
}

// Get returns a pointer to a Blower for a given BlowerId
func (b BlowerID) Get() *Blower {
	for k := range c.Blowers {
		if c.Blowers[k].ID == b {
			return &c.Blowers[k]
		}
	}
	return nil
}

// canEnable checks to see if a blower can be started given current global system state
// Blowers can run if the loop is not being actively serviced by a pump, but this is unusual
// blowers must start before pumps
// other than ensuring the system is on and the blower has not recently been stopped, there are no restrictions on starting a blower
func (b BlowerID) canEnable() error {
	if c.ControlMode == ControlOff {
		err := fmt.Errorf("system off, not starting blower")
		return err
	}

	blower := b.Get()
	if !blower.LastStopTime.Before(time.Now().Add(blowerMinTimeBetweenRuns)) {
		err := fmt.Errorf("blower recently stopped, in hold-down state")
		return err
	}

	return nil
}

// writeToStore writes the current running state to the long-term storage
// called at process shutdown and when a blower is stopped
func (b *Blower) writeToStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("blower-%d.json", b.ID))

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	j, err := json.Marshal(*b)
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

// readFromStore reads the last stored state into memory
// called at process startup
func (b *Blower) readFromStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("blower-%d.json", b.ID))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	var in Blower
	if err = json.Unmarshal(data, &in); err != nil {
		log.Error(err.Error())
		return err
	}

	// config file wins for id/name/hot/loop, only restore values which don't belong in the config
	b.LastStartTime = in.LastStartTime
	b.LastStopTime = in.LastStopTime
	b.Runtime = in.Runtime
	b.FilterTime = in.FilterTime

	return nil
}

// Start verifies that a blower can be enabled and sends the command to the MQTT subsystem, we do not record the activation of the blower, that takes place when the relay-module confirms that the blower has actually started
func (b BlowerID) Start(duration time.Duration, source string) error {
	if err := b.canEnable(); err != nil {
		log.Error("cannot enable blower", "id", b, "err", err.Error())
		return err
	}

	if duration < MinBlowerRunTime {
		err := fmt.Errorf("duration shorter than minimum: requested %.2f min %.2f", duration.Minutes(), MinBlowerRunTime.Minutes())
		return err
	}

	if duration > MaxBlowerRunTime {
		err := fmt.Errorf("duration longer than maximum: requested %.2f min %.2f", duration.Minutes(), MaxBlowerRunTime.Minutes())
		return err
	}

	cc := MQTTRequest{
		DeviceID: b,
		Command: Command{
			TargetState: true,
			RunTime:     duration,
			Source:      source,
		},
	}
	cmdChan <- cc
	return nil
}

// Stop sends the shutdown command to the MQTT subsystem to be sent to the relay-module
// if there are any pumps on the connected loops running, and no other blowers on the loop are running, the pump is stopped
func (b BlowerID) Stop(source string) {
	// if we are the last active blower on the loop, ensure that the pump is shut down
	last := true
	blower := b.Get()
	for k := range c.Blowers {
		// skip self
		if c.Blowers[k].ID == b {
			continue
		}
		// if anything else on the same hot loop is running, skip
		if c.SystemMode == SystemModeHeat && c.Blowers[k].HotLoop == blower.HotLoop && c.Blowers[k].Running {
			last = false
			break
		}
		// if anything else on the same cold loop is running, skip
		if c.SystemMode == SystemModeCool && c.Blowers[k].ColdLoop == blower.ColdLoop && c.Blowers[k].Running {
			last = false
			break
		}
	}

	if last {
		log.Debug("last running blower on loop, getting pump")
		pl := PumpID(0)
		if c.SystemMode == SystemModeHeat {
			pl = c.getPumpFromLoop(blower.HotLoop)
		} else {
			pl = c.getPumpFromLoop(blower.ColdLoop)
		}
		pump := pl.Get()
		log.Debug("pump", "id", pl, "pump", pump)
		if pump.Running {
			pump.ID.Stop("internal")
		}
	}

	cc := MQTTRequest{
		DeviceID: b,
		Command: Command{
			TargetState: false,
			RunTime:     0,
			Source:      source,
		},
	}
	cmdChan <- cc
}

func (b BlowerID) getPump(sm SystemModeT) PumpID {
	blower := b.Get()
	return blower.getPump(sm)
}

func (b *Blower) getPump(sm SystemModeT) PumpID {
	loopID := b.HotLoop
	if sm == SystemModeCool {
		loopID = b.ColdLoop
	}
	for k := range c.Pumps {
		if c.Pumps[k].Loop == loopID {
			return c.Pumps[k].ID
		}
	}
	log.Warn("unknown pump for blower", "blowerID", b.ID)
	return PumpID(0)
}
