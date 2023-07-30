package cmd

import (
	"github.com/lollipopkit/server_box_monitor/model"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:    "conf",
		Aliases: []string{"c"},
		Usage:   "Config file related commands",
		Subcommands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initialize config file",
				Action:  handleConfInit,
			},
		},
	})
}

func handleConfInit(c *cli.Context) error {
	return model.InitConfig()
}
