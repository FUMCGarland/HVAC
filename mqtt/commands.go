package hvacmqtt

import (
	"encoding/json"
	"fmt"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

// SendPumpTargetState sends a command to a pump
func SendPumpTargetState(p hvac.PumpID, cmd *hvac.PumpCommand) error {
	if cmd.TargetState {
		if err := p.CanEnable(); err != nil {
			log.Error("cannot enable pump", "id", p, "cmd", cmd, "err", err.Error())
			return err
		}
	}

	topic := fmt.Sprintf("%s/pumps/%d/targetstate", root, p)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Info("Pump TargetState", "pump", p, "target", cmd.TargetState, "topic", topic)

	return inline.Publish(topic, jcmd, false, 0)
}

// SendBlowerTargetState sends a command to a pump
func SendBlowerTargetState(b hvac.BlowerID, cmd *hvac.BlowerCommand) error {
	c := hvac.GetConfig()
	if cmd.TargetState {
		if err := b.CanEnable(); err != nil {
			log.Error("cannot enable blower", "id", b, "cmd", cmd, "err", err.Error())
			return err
		}
	} else {
		// TODO: this logic needs to be moved to hvac.BlowerID.Stop()
		// if we are the last active blower on the loop, ensure that the pump is shut down
		last := true
		blower := b.Get()
		for k := range c.Blowers {
			// skip self
			if c.Blowers[k].ID == b {
				continue
			}
			// if anything else on the same hot loop is running, skip
			if c.SystemMode == hvac.SystemModeHeat && c.Blowers[k].HotLoop == blower.HotLoop && c.Blowers[k].Running {
				last = false
				break
			}
			// if anything else on the same cold loop is running, skip
			if c.SystemMode == hvac.SystemModeCool && c.Blowers[k].ColdLoop == blower.ColdLoop && c.Blowers[k].Running {
				last = false
				break
			}
		}

		if last {
			pl := hvac.PumpID(0)
			if c.SystemMode == hvac.SystemModeHeat {
				pl = c.GetPumpFromLoop(blower.HotLoop)
			} else {
				pl = c.GetPumpFromLoop(blower.ColdLoop)
			}
			pump := pl.Get()
			if pump.Running {
				pump.ID.Stop("internal")
			}
		}
	}

	topic := fmt.Sprintf("%s/blowers/%d/targetstate", root, b)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Info("Blower TargetState", "blower", b, "target", cmd.TargetState, "topic", topic)

	return inline.Publish(topic, jcmd, false, 0)
}
