package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

type ChillerID uint8

var chillerLockout bool

type Chiller struct {
	ID               ChillerID
	Name             string
	Loops            []LoopID
	Runtime          time.Duration
	Running          bool
	CurrentStartTime time.Time
	LastStartTime    time.Time
	LastStopTime     time.Time
}

func (p ChillerID) Get() *Chiller {
	for k := range c.Chillers {
		if c.Chillers[k].ID == p {
			return &c.Chillers[k]
		}
	}
	return nil
}

// canEnable
// (1) the SystemMode must be cooling
// (2) a cool pump must be running (which implies a blower running)
// (3) don't fast-cycle the chillers, if you stop it, leave it stopped for at least 5(?) minutes
// (4) if anything is below 60degF, do not enable the chiller lest it freeze out
func (ch ChillerID) canEnable() error {
	if c.ControlMode == ControlOff {
		err := fmt.Errorf("system off, not starting chiller")
		return err
	}

	chiller := ch.Get()
	if c.SystemMode != SystemModeCool {
		err := fmt.Errorf("cannot enable chiller if not in cooling")
		return err
	}

	pumpRunning := false
	for k := range c.Pumps {
		if c.Pumps[k].SystemMode == SystemModeCool && c.Pumps[k].Running {
			pumpRunning = true
		}
	}
	if !pumpRunning {
		err := fmt.Errorf("cannot enable chiller if no pumps are running")
		return err
	}

	if !chiller.LastStopTime.Before(time.Now().Add(chillerMinTimeBetweenRuns)) {
		err := fmt.Errorf("chiller recently stopped, in hold-down state")
		return err
	}

	if chillerLockout {
		chillerReset := true
		// TODO make sure the temp reports are recent
		for k := range c.Rooms {
			if c.Rooms[k].Temperature != 0 && c.Rooms[k].Temperature < chillerRecoveryTemp {
				// a room below the reset temp, do not reset
				chillerReset = false
				break
			}
		}

		if chillerReset {
			log.Info("all rooms above recovery temp, unlocking chiller")
			chillerLockout = false
		}
	}

	if chillerLockout {
		err := fmt.Errorf("chiller locked out, not enabling")
		log.Error(err.Error())
		return err
	}

	return nil
}

func (ch *Chiller) writeToStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("chiller-%d.json", ch.ID))

	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	j, err := json.Marshal(*ch)
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

func (ch *Chiller) readFromStore() error {
	path := path.Join(c.StateStore, fmt.Sprintf("chiller-%d.json", ch.ID))

	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	var in Chiller
	if err = json.Unmarshal(data, &in); err != nil {
		log.Error(err.Error())
		return err
	}

	// config file wins for id/name/hot/loop, only restore values which don't belong in the config
	ch.LastStartTime = in.LastStartTime
	ch.LastStopTime = in.LastStopTime
	ch.Runtime = in.Runtime

	return nil
}

func (ch ChillerID) Start(duration time.Duration, source string) error {
	if err := ch.canEnable(); err != nil {
		log.Error("cannot enable ch", "id", ch, "err", err.Error())
		return err
	}

	if duration < MinChillerRunTime {
		err := fmt.Errorf("duration shorter than minimum: requested %.2f min %.2f", duration.Minutes(), MinChillerRunTime.Minutes())
		return err
	}

	if duration > MaxChillerRunTime {
		err := fmt.Errorf("duration longer than maximum: requested %.2f min %.2f", duration.Minutes(), MaxChillerRunTime.Minutes())
		return err
	}

	cc := MQTTRequest{
		Device: ch,
		Command: Command{
			TargetState: true,
			RunTime:     duration,
			Source:      source,
		},
	}
	cmdChan <- cc
	return nil
}

func (ch ChillerID) Stop(source string) {
	cc := MQTTRequest{
		Device: ch,
		Command: Command{
			TargetState: false,
			RunTime:     0,
			Source:      source,
		},
	}
	cmdChan <- cc
}

func (c *Config) GetChillerFromLoop(id LoopID) ChillerID {
	for k := range c.Chillers {
		for j := range c.Chillers[k].Loops {
			if c.Chillers[k].Loops[j] == id {
				return c.Chillers[k].ID
			}
		}
	}
	return ChillerID(0)
}

func (ch ChillerID) PumpsRunning() bool {
	chiller := ch.Get()
	for _, l := range chiller.Loops {
		p := c.getPumpFromLoop(l)
		pump := p.Get()
		if pump.Running {
			return true
		}
	}
	return false
}
