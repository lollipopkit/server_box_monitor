package runner

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lollipopkit/server_box_monitor/model"
	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	pushArgs     = []*model.PushFormatArgs{}
	pushArgsLock = new(sync.RWMutex)
)

func init() {
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

func Start() {
	go Run()
	go Push()
	// 阻塞主线程
	select {}
}

func Run() {
	for {
		err := model.ReadAppConfig()
		if err != nil {
			time.Sleep(model.DefaultappConfig.GetRunInterval())
			continue
		}

		err = model.RefreshStatus()
		status := model.GetStatus()
		if err != nil {
			utils.Warn("[STATUS] Get status error: %v", err)
			goto SLEEP
		}

		for _, rule := range model.Config.Rules {
			notify, arg, err := rule.ShouldNotify(status)
			if err != nil {
				if !strings.Contains(err.Error(), "not ready") {
					utils.Warn("[RULE] %s error: %v", rule.Id(), err)
				}
			}

			if notify && arg != nil {
				pushArgsLock.Lock()
				pushArgs = append(pushArgs, arg)
				pushArgsLock.Unlock()
			}
		}

		// utils.Info("[STATUS] refreshed, %d to push", len(pushArgs))
	SLEEP:
		time.Sleep(model.Config.GetRunInterval())
		continue
	}
}

func Push() {
	for {
		err := model.ReadAppConfig()
		if err != nil {
			time.Sleep(model.DefaultappConfig.GetRunInterval())
			continue
		}

		if len(pushArgs) == 0 {
			time.Sleep(model.Config.GetPushInterval())
			continue
		}

		for _, push := range model.Config.Pushes {
			pushArgsLock.RLock()
			err := push.Push(pushArgs)
			pushArgsLock.RUnlock()
			if err != nil {
				utils.Warn("[PUSH] %s error: %v", push.Id(), err)
				continue
			}
			utils.Success("[PUSH] %s success", push.Id())
		}
		pushArgsLock.Lock()
		pushArgs = []*model.PushFormatArgs{}
		pushArgsLock.Unlock()

		time.Sleep(model.Config.GetPushInterval())
	}
}
