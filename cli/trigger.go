package cli

import (
	"log"

	"github.com/codegangsta/cli"
)

func trigger(c *cli.Context) {
	// Input validation
	inventoryBase64Str := c.String(InventoryDataKey)

	configBase64Str := c.String(ConfigDataKey)
	if configBase64Str == "" {
		log.Fatal("Missing required config data")
	}

	triggerPattern := c.String(TriggerPatternNameKey)
	if triggerPattern == "" {
		log.Fatal("Missing required Trigger Pattern parameter")
	}

	if _, _, err := performRunOrTrigger(inventoryBase64Str, configBase64Str, triggerPattern, true); err != nil {
		log.Fatalf("Failed to run, error: %s", err)
	}
}
