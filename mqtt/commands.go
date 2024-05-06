package hvacmqtt

import (
	"encoding/json"
	"fmt"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

// SendPumpTargetState sends a command to a pump
func sendPumpTargetState(p hvac.PumpID, cmd *hvac.PumpCommand) error {
	topic := fmt.Sprintf("%s/pumps/%d/targetstate", root, p)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Info("controller: Pump TargetState", "pump", p, "target", cmd.TargetState, "topic", topic)
	return inline.Publish(topic, jcmd, false, 0)
}

// SendBlowerTargetState sends a command to a pump
func sendBlowerTargetState(b hvac.BlowerID, cmd *hvac.BlowerCommand) error {
	topic := fmt.Sprintf("%s/blowers/%d/targetstate", root, b)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Info("controller: blower TargetState", "blower", b, "target", cmd.TargetState, "topic", topic)
	return inline.Publish(topic, jcmd, false, 0)
}

// SendChillerTargetState sends a command to a pump
func sendChillerTargetState(ch hvac.ChillerID, cmd *hvac.ChillerCommand) error {
	topic := fmt.Sprintf("%s/chillers/%d/targetstate", root, ch)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Info("controller: Chiller TargetState", "chiller", ch, "target", cmd.TargetState, "topic", topic)
	return inline.Publish(topic, jcmd, false, 0)
}
