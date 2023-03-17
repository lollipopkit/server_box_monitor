package cmd

import (
	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/urfave/cli/v2"
)

func init() {
	cmds = append(cmds, &cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Config something",
		Subcommands: []*cli.Command{
			{
				Name:      "alias",
				Aliases:   []string{"a"},
				Usage:     "Config alias `dart -> fvm dart` and `flutter -> fvm flutter`",
				Action:    handlePlaceHolder,
				ArgsUsage: "[alias]",
				UsageText: res.APP_NAME + " config alias [alias]",
			},
			{
				Name:      "use-mirror",
				Aliases:   []string{"um"},
				Usage:     "config use mirror or not",
				Action:    handlePlaceHolder,
				UsageText: res.APP_NAME + " config use-mirror [ true | false ]",
			},
		},
	})
}
