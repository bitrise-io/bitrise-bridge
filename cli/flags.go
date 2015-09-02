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

	// TriggerPatternNameKey ...
	TriggerPatternNameKey = "pattern"
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
		Usage: "inventory/secrets data (~ content of .bitrise.secrets.yml)",
	}

	flConfigData = cli.StringFlag{
		Name:  ConfigDataKey + ", " + configDataKeyShort,
		Usage: "config data (~ content of bitrise.yml)",
	}

	flWorkflowName = cli.StringFlag{
		Name:  WorkflowNameKey + ", " + workflowNameKeyShort,
		Usage: "workflow to pass to bitrise",
	}

	flTriggerPattern = cli.StringFlag{
		Name:  TriggerPatternNameKey,
		Usage: "trigger pattern to pass to bitrise",
	}
)
