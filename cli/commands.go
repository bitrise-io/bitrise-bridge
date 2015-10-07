package cli

import "github.com/codegangsta/cli"

const (
	// LogLevelKey ...
	LogLevelKey      = "loglevel"
	logLevelKeyShort = "l"
	// LogLevelEnvKey ...
	LogLevelEnvKey string = "LOGLEVEL"
	// CommandHostKey ...
	CommandHostKey = "command-host"

	// InventoryDataKey ...
	InventoryDataKey = "inventory"
	// InventoryDataBase64Key ...
	InventoryDataBase64Key = "inventory-base64"
	inventoryDataKeyShort  = "i"

	// ConfigDataKey ...
	ConfigDataKey = "config"
	// ConfigDataBase64Key ...
	ConfigDataBase64Key = "config-base64"
	configDataKeyShort  = "c"

	// WorkflowNameKey ...
	WorkflowNameKey      = "workflow"
	workflowNameKeyShort = "w"

	// TriggerPatternNameKey ...
	TriggerPatternNameKey = "pattern"

	// DockerImageIDNameKey ...
	DockerImageIDNameKey = "docker-image-id"

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
	flCommandHost = cli.StringFlag{
		Name:  CommandHostKey,
		Usage: "Command host. (none, cmd-bridge [default, for compatibility], docker)",
	}
	flDockerImageID = cli.StringFlag{
		Name:  DockerImageIDNameKey,
		Usage: "docker image ID to use - only applies to command-host=docker",
	}
	flags = []cli.Flag{
		flLogLevel,
		flCommandHost,
		flDockerImageID,
	}

	// command flags
	flInventoryData = cli.StringFlag{
		Name:  InventoryDataKey + ", " + inventoryDataKeyShort + ", " + InventoryDataBase64Key,
		Usage: "inventory/secrets data in base64 (~ content of .bitrise.secrets.yml)",
	}

	flConfigData = cli.StringFlag{
		Name:  ConfigDataKey + ", " + configDataKeyShort + ", " + ConfigDataBase64Key,
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

var (
	commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "calls: bitrise run",
			Action:  run,
			Flags: []cli.Flag{
				flInventoryData,
				flConfigData,
				flWorkflowName,
				flWorkdirPath,
			},
		},
		{
			Name:   "trigger",
			Usage:  "calls: bitrise trigger",
			Action: trigger,
			Flags: []cli.Flag{
				flInventoryData,
				flConfigData,
				flTriggerPattern,
				flWorkdirPath,
			},
		},
	}
)
