package runner

import (
	"os"
	"time"

	"github.com/lollipopkit/server_box_monitor/model"
	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

func init() {
	// Check script file
	if utils.Exist(res.ServerBoxShellPath) {
		utils.Info("[INIT] Script has been installed.")
		return
	}

	// Write script file
	scriptBytes, err := res.Files.ReadFile(res.ServerBoxShellFileName)
	if err != nil {
		utils.Error("[INIT] Read embed file error: %v", err)
		panic(err)
	}
	err = os.WriteFile(res.ServerBoxShellPath, scriptBytes, 0755)
	if err != nil {
		utils.Error("[INIT] Write script file error: %v", err)
		panic(err)
	}
}

func Run() {
	for {
		appConfig, err := model.ReadAppConfig()
		if err != nil {
			utils.Error("[CONFIG] Read file error: %v", err)
			time.Sleep(model.DefaultappConfig.GetRunInterval())
			continue
		}

		args := []*model.PushFormatArgs{}

		for _, rule := range appConfig.Rules {
			status, err := model.GetStatus()
			if err != nil {
				utils.Warn("[STATUS] Get status error: %v", err)
				time.Sleep(appConfig.GetRunInterval())
				continue
			}
			notify, arg, err := rule.ShouldNotify(status)
			if err != nil {
				utils.Warn("[RULE] check error: %v", err)
				time.Sleep(appConfig.GetRunInterval())
				continue
			}

			if notify {
				args = append(args, arg)
			}
		}

		for _, push := range appConfig.Pushes {
			err := push.Push(args)
			if err != nil {
				utils.Warn("[PUSH] Push error: %v", err)
			}
		}
	}
}