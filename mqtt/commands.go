package hvacmqtt

import (
	"encoding/json"
	"fmt"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

// sendTargetState sends a generic command to a generic device
func sendTargetState(did hvac.DeviceID, cmd *hvac.Command) error {
	subtopic := hvac.PumpsTopic
	switch did.(type) {
	case hvac.PumpID:
		subtopic = hvac.PumpsTopic
	case hvac.BlowerID:
		subtopic = hvac.BlowersTopic
	case hvac.ChillerID:
		subtopic = hvac.ChillersTopic
	}

	topic := fmt.Sprintf("%s/%s/%d/%s", root, subtopic, did, hvac.TargetStateEndpoint)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Debug("controller sending targetState", "deviceID", did, "target", cmd.TargetState, "topic", topic)
	return inline.Publish(topic, jcmd, false, hvac.QoS)
}
