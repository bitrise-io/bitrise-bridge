package cli

import "github.com/codegangsta/cli"

const (
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"
	// LogLevelEnvKey ...
	LogLevelEnvKey string = "LOGLEVEL"

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

	// WorkdirPathKey ...
	WorkdirPathKey = "workdir"
	// WorkdirPathEnvKey ...
	WorkdirPathEnvKey = "BITRISE_BRIDGE_WORKDIR"
)

var (
	// App flags
	flLogLevel = cli.StringFlag{
		Name:   LogLevelKey + ", " + logLevelKeyShort,
		Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
		EnvVar: LogLevelEnvKey,
	}
	flags = []cli.Flag{
		flLogLevel,
	}

	flInventoryData = cli.StringFlag{
		Name:  InventoryDataKey + ", " + inventoryDataKeyShort,
		Usage: "inventory/secrets data in base64 (~ content of .bitrise.secrets.yml)",
	}

	flConfigData = cli.StringFlag{
		Name:  ConfigDataKey + ", " + configDataKeyShort,
		Usage: "config data in base64 (~ content of bitrise.yml)",
	}

	flWorkflowName = cli.StringFlag{
		Name:  WorkflowNameKey + ", " + workflowNameKeyShort,
		Usage: "workflow to pass to bitrise",
	}

	flTriggerPattern = cli.StringFlag{
		Name:  TriggerPatternNameKey,
		Usage: "trigger pattern to pass to bitrise",
	}

	flWorkdirPath = cli.StringFlag{
		Name:   WorkdirPathKey,
		Usage:  "set this to bitrise as the working directory",
		EnvVar: WorkdirPathEnvKey,
	}
)
