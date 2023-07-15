package main

import (
	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/server_box_monitor/cmd"
)

func main() {
	cmd.Run()

	log.Setup(log.Config{
		PrintTime: true,
	})
}
