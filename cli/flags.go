package cli

import "github.com/codegangsta/cli"

const (
	// ConfigDataKey ...
	ConfigDataKey      = "config-data"
	configDataKeyShort = "c"

	// InventoryDataKey ...
	InventoryDataKey      = "inventory-data"
	inventoryDataKeyShort = "i"
)

var (
	flags = []cli.Flag{}

	flConfigData = cli.StringFlag{
		Name:  ConfigDataKey + ", " + configDataKeyShort,
		Usage: ".",
	}

	flInventoryData = cli.StringFlag{
		Name:  InventoryDataKey + ", " + inventoryDataKeyShort,
		Usage: ".",
	}
)
