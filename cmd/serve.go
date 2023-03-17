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
	})
}

func handleServe(ctx *cli.Context) error {
	runner.Start()
	return nil
}
