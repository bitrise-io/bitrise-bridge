package cli

import "github.com/codegangsta/cli"

var (
	commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   ".",
			Action:  run,
			Flags: []cli.Flag{
				flInventoryData,
				flConfigData,
				flWorkflowName,
			},
		},
	}
)
