package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-bridge/bridge"
	"github.com/bitrise-io/bitrise-bridge/config"
	"github.com/bitrise-io/go-utils/parseutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/codegangsta/cli"
)

var (
	// BridgeConfigs ...
	BridgeConfigs config.Model
)

func initLogFormatter() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})
}

func before(c *cli.Context) error {
	initLogFormatter()

	// Log level
	logLevelStr := c.String(LogLevelKey)
	if logLevelStr == "" {
		logLevelStr = "info"
	}

	level, err := log.ParseLevel(logLevelStr)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	// Command Host
	commandHostStr := c.String(CommandHostKey)
	if commandHostStr != "" {
		switch commandHostStr {
		case bridge.CommandHostIDCmdBridge, bridge.CommandHostIDNone, bridge.CommandHostIDDocker:
			CommandHostID = commandHostStr
		default:
			log.Fatalf("Invalid / not supported command-host specified: %s", commandHostStr)
		}
	}
	log.Debugf("Command host: %s", CommandHostID)

	// Load config file, if any
	configFilePath := config.DefaultConfigFilePath()
	if exist, err := pathutil.IsPathExists(configFilePath); err != nil {
		log.Fatalf("Failed to check config file: %s", err)
	} else if exist {
		log.Debugf("Using config found at path: %s", configFilePath)
		conf, err := config.LoadConfigFromFile(configFilePath)
		if err != nil {
			log.Fatalf("Failed to read config file: %s", err)
		}
		BridgeConfigs = conf
	}

	// Command Host Args
	if CommandHostID == bridge.CommandHostIDDocker {
		commandHostDockerImage := c.String(DockerImageIDNameKey)
		if commandHostDockerImage != "" {
			BridgeConfigs.Docker.Image = commandHostDockerImage
		}
		commandHostDockerAllowAccessToDockerInContainer := c.String(DockerAllowAccessToDockerInContainer)
		if commandHostDockerAllowAccessToDockerInContainer != "" {
			val, err := parseutil.ParseBool(commandHostDockerAllowAccessToDockerInContainer)
			if err != nil {
				log.Warnf("Invalid parameter 'docker-allow-access-to-docker-in-container': %s", commandHostDockerAllowAccessToDockerInContainer)
				log.Warn("=> Ignoring the parameter")
			} else {
				BridgeConfigs.Docker.IsAllowAccessToDockerInContainer = val
			}
		}
	}
	log.Debugf("Configs: %#v", BridgeConfigs)

	return nil
}

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

// Run ...
func Run() {
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = ""
	app.Version = "0.9.7"

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Finished with Error:", err)
	}
}
