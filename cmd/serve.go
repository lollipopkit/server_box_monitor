package cmd

import (
	"github.com/lollipopkit/server_box_monitor/runner"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "Run monitor",
		Action:  handleServe,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addr",
				Aliases: []string{"a"},
				Usage:   "Listen address",
				Value:   "0.0.0.0:3770",
				EnvVars: []string{"SBM_ADDR"},
			},
		},
	})
}

func handleServe(ctx *cli.Context) error {
	runner.Start(ctx.String("addr"))
	return nil
}
