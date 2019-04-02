package cli

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-io/bitrise-bridge/bridge"
	"github.com/codegangsta/cli"
)

func run(c *cli.Context) error {
	// Input validation
	inventoryBase64Str := c.String(InventoryDataKey)

	configBase64Str := c.String(ConfigDataKey)
	if configBase64Str == "" {
		return errors.New("Missing required config data")
	}

	runParamJSONBase64 := c.String(JSONParamsBase64Key)

	workflowName := c.String(WorkflowNameKey)
	if workflowName == "" && runParamJSONBase64 == "" {
		return errors.New("Missing required workflow name or JSON (base64) run params (at least one is required)")
	}

	if err := bridge.PerformRunOrTrigger(CommandHostID, BridgeConfigs, inventoryBase64Str, configBase64Str, runParamJSONBase64, workflowName, false, c.String(WorkdirPathKey)); err != nil {
		if !errorutil.IsExitStatusError(err) {
			log.Errorf("Failed to run, error: %s", err)
		}
		return errors.New("")
	}

	return nil
}
