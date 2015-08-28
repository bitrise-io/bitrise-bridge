package cli

import (
	"github.com/bitrise-io/bitrise/bitrise"
	"github.com/codegangsta/cli"
)

const (
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"

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
	// App flags
	flLogLevel = cli.StringFlag{
		Name:   LogLevelKey + ", " + logLevelKeyShort,
		Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
		EnvVar: bitrise.LogLevelEnvKey,
	}
	flags = []cli.Flag{
		flLogLevel,
	}

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
