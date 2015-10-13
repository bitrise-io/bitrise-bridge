package cli

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-bridge/bridge"
	"github.com/codegangsta/cli"
)

var (
	// CommandHostArgs ...
	CommandHostArgs = map[string]string{}
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

	// Command Host Args
	if CommandHostID == bridge.CommandHostIDDocker {
		commandHostDockerImage := c.String(DockerImageIDNameKey)
		if commandHostDockerImage != "" {
			CommandHostArgs["docker-image-id"] = commandHostDockerImage
		}
	}
	log.Debugf("CommandHostArgs: %#v", CommandHostArgs)

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
	app.Version = "0.9.5"

	app.Author = ""
	app.Email = ""

	app.Before = before

	app.Flags = flags
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Finished with Error:", err)
	}
}
