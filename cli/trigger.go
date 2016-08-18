package cli

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-tools/bitrise-bridge/bridge"
	"github.com/codegangsta/cli"
)

func trigger(c *cli.Context) error {
	// Input validation
	inventoryBase64Str := c.String(InventoryDataKey)

	configBase64Str := c.String(ConfigDataKey)
	if configBase64Str == "" {
		return errors.New("Missing required config data")
	}

	triggerPattern := c.String(TriggerPatternNameKey)
	if triggerPattern == "" {
		return errors.New("Missing required Trigger Pattern parameter")
	}

	if err := bridge.PerformRunOrTrigger(CommandHostID, BridgeConfigs, inventoryBase64Str, configBase64Str, triggerPattern, true, c.String(WorkdirPathKey)); err != nil {
		if !errorutil.IsExitStatusError(err) {
			log.Errorf("Failed to trigger, error: %s", err)
		}
		return errors.New("")
	}

	return nil
}
