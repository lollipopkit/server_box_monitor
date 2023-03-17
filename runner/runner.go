package runner

import (
	"os"
	"time"

	"github.com/lollipopkit/server_box_monitor/model"
	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

func init() {
	if utils.Exist(res.ServerBoxShellPath) {
		return
	}

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
			time.Sleep(model.DefaultappConfig.GetRunInterval())
			continue
		}

		args := []*model.PushFormatArgs{}
		status, err := model.GetStatus()
		if err != nil {
			utils.Warn("[STATUS] Get status error: %v", err)
			goto SLEEP
		}

		for _, rule := range appConfig.Rules {
			notify, arg, err := rule.ShouldNotify(status)
			if err != nil {
				utils.Warn("[RULE] %s error: %v", rule.Id(), err)
			}

			if notify && arg != nil {
				args = append(args, arg)
			}
		}

		if len(args) == 0 {
			goto SLEEP
		}

		for _, push := range appConfig.Pushes {
			err := push.Push(args)
			if err != nil {
				utils.Warn("[PUSH] %s error: %v", push.Id(), err)
				continue
			}
			utils.Success("[PUSH] %s success", push.Id())
		}

	SLEEP:
		time.Sleep(appConfig.GetRunInterval())
		continue
	}
}
