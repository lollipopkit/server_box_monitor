package cmd

import (
	"github.com/lollipopkit/server_box_monitor/model"
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
			&cli.StringFlag{
				Name:    "crt",
				Aliases: []string{"c"},
				Usage:   "TLS certificate file path",
				EnvVars: []string{"SBM_TLS_CRT"},
			},
			&cli.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Usage:   "TLS key file path",
				EnvVars: []string{"SBM_TLS_KEY"},
			},
		},
	})
}

func handleServe(ctx *cli.Context) error {
	webConfig := &model.WebConfig{
		Addr: ctx.String("addr"),
		Cert: ctx.String("crt"),
		Key:  ctx.String("key"),
	}
	runner.Start(webConfig)
	return nil
}
