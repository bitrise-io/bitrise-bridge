package cli

import "github.com/codegangsta/cli"

const (
	// InventoryDataKey ...
	InventoryDataKey      = "inventory"
	inventoryDataKeyShort = "i"

	// ConfigDataKey ...
	ConfigDataKey      = "config"
	configDataKeyShort = "c"

	// WorkflowNameKey ...
	WorkflowNameKey      = "workflow"
	workflowNameKeyShort = "w"
)

var (
	flags = []cli.Flag{}

	flInventoryData = cli.StringFlag{
		Name:  InventoryDataKey + ", " + inventoryDataKeyShort,
		Usage: ".",
	}

	flConfigData = cli.StringFlag{
		Name:  ConfigDataKey + ", " + configDataKeyShort,
		Usage: ".",
	}

	flWorkflowName = cli.StringFlag{
		Name:  WorkflowNameKey + ", " + workflowNameKeyShort,
		Usage: ".",
	}
)
