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
	pushPairs     = []*model.PushPair{}
	pushPairsLock = new(sync.RWMutex)
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
	// 阻塞主线程
	select {}
}

func Run() {
	err := model.ReadAppConfig()
	if err != nil {
		utils.Error("[CONFIG] Read app config error: %v", err)
		panic(err)
	}
	
	for {
		err = model.RefreshStatus()
		status := model.GetStatus()
		if err != nil {
			utils.Warn("[STATUS] Get status error: %v", err)
			goto SLEEP
		}

		for _, rule := range model.Config.Rules {
			notify, pushPair, err := rule.ShouldNotify(status)
			if err != nil {
				if !strings.Contains(err.Error(), "not ready") {
					utils.Warn("[RULE] %s error: %v", rule.Id(), err)
				}
			}

			if notify && pushPair != nil {
				pushPairsLock.Lock()
				pushPairs = append(pushPairs, pushPair)
				pushPairsLock.Unlock()
			}
		}

		if len(pushPairs) == 0 {
			goto SLEEP
		}

		utils.Info("[STATUS] refreshed, %d to push", len(pushPairs))

		pushPairsLock.RLock()
		for _, push := range model.Config.Pushes {
			err := push.Push(pushPairs)
			if err != nil {
				utils.Warn("[PUSH] %s error: %v", push.Id(), err)
				continue
			}
			utils.Success("[PUSH] %s success", push.Id())
		}
		pushPairsLock.RUnlock()

		pushPairsLock.Lock()
		pushPairs = []*model.PushPair{}
		pushPairsLock.Unlock()
	SLEEP:
		time.Sleep(model.GetInterval())
	}
}
