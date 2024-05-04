package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

type BlowerID uint8

type Blower struct {
	ID               BlowerID
	Name             string
	HotLoop          LoopID
	ColdLoop         LoopID
	Zone             ZoneID
	Running          bool
	Runtime          time.Duration
	FilterTime       uint64
	CurrentStartTime time.Time
	LastStartTime    time.Time
	LastStopTime     time.Time
}

func (b BlowerID) Get() *Blower {
	for k := range c.Blowers {
		if c.Blowers[k].ID == b {
			return &c.Blowers[k]
		}
	}
	return nil
}

func (b BlowerID) CanEnable() error {
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

func (b BlowerID) Start(duration time.Duration, source string) error {
	if err := b.CanEnable(); err != nil {
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
		Device: b,
		Command: Command{
			TargetState: true,
			RunTime:     duration,
			Source:      source,
		},
	}
	cmdChan <- cc
	return nil
}

func (b BlowerID) Stop(source string) {
	cc := MQTTRequest{
		Device: b,
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
	log.Info("unknown pump for blower", "blowerID", b.ID)
	return PumpID(0)
}
