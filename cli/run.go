package cli

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/errorutil"
	"github.com/bitrise-tools/bitrise-bridge/bridge"
	"github.com/codegangsta/cli"
)

func run(c *cli.Context) {
	// Input validation
	inventoryBase64Str := c.String(InventoryDataKey)

	configBase64Str := c.String(ConfigDataKey)
	if configBase64Str == "" {
		log.Fatal("Missing required config data")
	}

	workflowName := c.String(WorkflowNameKey)
	if workflowName == "" {
		log.Fatal("Missing required workflow name")
	}

	if err := bridge.PerformRunOrTrigger(CommandHostID, BridgeConfigs, inventoryBase64Str, configBase64Str, workflowName, false, c.String(WorkdirPathKey)); err != nil {
		if !errorutil.IsExitStatusError(err) {
			log.Errorf("Failed to run, error: %s", err)
		}
		os.Exit(1)
	}
}
