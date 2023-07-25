package cmd

import (
	// "strings"

	// "github.com/lollipopkit/gommon/res"
	// "github.com/lollipopkit/gommon/term"
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
			// {
			// 	Name:    "edit",
			// 	Aliases: []string{"e"},
			// 	Usage:   "Edit config file",
			// 	Action:  handleConfEdit,
			// },
		},
	})
}

func handleConfInit(c *cli.Context) error {
	return model.InitConfig()
}

// func handleConfEdit(c *cli.Context) error {
// 	if err := model.ReadAppConfig(); err != nil {
// 		return err
// 	}

// 	typeOptions := []string{"interval", "rate", "name", "rules", "pushes", "exit"}
// 	opOptions := []string{"add", "remove", "edit", "exit"}

// 	for {
// 		ruleIds := []string{}
// 		for _, rule := range model.Config.Rules {
// 			ruleIds = append(ruleIds, rule.Id())
// 		}
// 		pushNames := []string{}
// 		for _, push := range model.Config.Pushes {
// 			pushNames = append(pushNames, push.Name)
// 		}
// 		var buf strings.Builder
// 		buf.WriteString(res.GREEN + "interval: " + res.NOCOLOR + model.Config.Interval + "\n")
// 		buf.WriteString(res.GREEN + "rate: " + res.NOCOLOR + model.Config.Rate + "\n")
// 		buf.WriteString(res.GREEN + "name: " + res.NOCOLOR + model.Config.Name + "\n")
// 		buf.WriteString(res.GREEN + "rules: " + res.NOCOLOR + strings.Join(ruleIds, " | ") + "\n")
// 		buf.WriteString(res.GREEN + "pushes: " + res.NOCOLOR + strings.Join(pushNames, " | ") + "\n")
// 		print(buf.String())

// 		op := term.Option("What to do?", opOptions, len(opOptions)-1)
// 		if op == 3 {
// 			break
// 		}
// 		question := "Which type to " + opOptions[op] + "?"
// 		typ := term.Option(question, typeOptions, len(typeOptions)-1)
// 		switch typ {
// 		case 0:
			
// 	}
// 	return nil
// }
