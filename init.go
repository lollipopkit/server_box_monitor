package main

import (
	"os"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

func BeforeStart() {
	// Create serverbox dir
	if !utils.Exist(res.ServerBoxDirPath) {
		err := os.MkdirAll(res.ServerBoxDirPath, 0755)
		if err != nil {
			utils.Error("Create dir error: %v", err)
			panic(err)
		}
	}
}