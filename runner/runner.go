package runner

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lollipopkit/gommon/term"
	"github.com/lollipopkit/server_box_monitor/model"
	"github.com/lollipopkit/server_box_monitor/res"
)

var (
	pushPairs     = []*model.PushPair{}
	pushPairsLock = new(sync.RWMutex)
)

func init() {
	scriptBytes, err := res.Files.ReadFile(res.ServerBoxShellFileName)
	if err != nil {
		term.Err("[INIT] Read embed file error: %v", err)
		panic(err)
	}
	err = os.WriteFile(res.ServerBoxShellPath, scriptBytes, 0755)
	if err != nil {
		term.Err("[INIT] Write script file error: %v", err)
		panic(err)
	}
}

func Start() {
	go run()
	// 阻塞主线程
	select {}
}

func run() {
	err := model.ReadAppConfig()
	if err != nil {
		term.Err("[CONFIG] Read app config error: %v", err)
		panic(err)
	}

	for range time.NewTicker(model.GetInterval()).C {
		err = model.RefreshStatus()
		status := model.Status
		if err != nil {
			term.Warn("[STATUS] Get status error: %v", err)
			continue
		}

		for _, rule := range model.Config.Rules {
			notify, pushPair, err := rule.ShouldNotify(status)
			if err != nil {
				if !strings.Contains(err.Error(), "not ready") {
					term.Warn("[RULE] %s error: %v", rule.Id(), err)
				}
			}

			if notify && pushPair != nil {
				pushPairsLock.Lock()
				pushPairs = append(pushPairs, pushPair)
				pushPairsLock.Unlock()
			}
		}

		if len(pushPairs) == 0 {
			continue
		}

		term.Info("[STATUS] refreshed, %d to push", len(pushPairs))

		pushPairsLock.RLock()
		for _, push := range model.Config.Pushes {
			err := push.Push(pushPairs)
			if err != nil {
				term.Warn("[PUSH] %s error: %v", push.Name, err)
				continue
			}
			term.Suc("[PUSH] %s success", push.Name)
		}
		pushPairsLock.RUnlock()

		pushPairsLock.Lock()
		pushPairs = []*model.PushPair{}
		pushPairsLock.Unlock()
	}
}
