package cli

import "github.com/codegangsta/cli"

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
