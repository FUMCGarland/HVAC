package hvacmqtt

import (
	"encoding/json"
	"fmt"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

// sendTargetState sends a generic command to a generic device
func sendTargetState(did hvac.DeviceID, cmd *hvac.Command) error {
	devtype := "pump"
	switch did.(type) {
	case hvac.PumpID:
		devtype = "pumps"
	case hvac.BlowerID:
		devtype = "blowers"
	case hvac.ChillerID:
		devtype = "chillers"
	}

	topic := fmt.Sprintf("%s/%s/%d/targetstate", root, devtype, did)
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		log.Error("unable to marshal command", "cmd", cmd, "topic", topic)
		return err
	}

	log.Info("controller sending targetState", "type", devtype, "deviceID", did, "target", cmd.TargetState, "topic", topic)
	return inline.Publish(topic, jcmd, false, 0)
}
