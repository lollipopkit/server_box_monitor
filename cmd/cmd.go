package cmd

import (
	"os"

	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/urfave/cli/v2"
)

var (
	cmds  = []*cli.Command{}
	flags = []cli.Flag{}
)

func Run() {
	app := &cli.App{
		Name:     res.APP_NAME,
		Usage:    "a Monitor runs on server",
		Version:  res.APP_VERSION,
		Commands: cmds,
		Flags:    flags,
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() == 0 {
				return cli.ShowAppHelp(ctx)
			}
			return nil
		},
		Suggest: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Err(err.Error())
	}
}
