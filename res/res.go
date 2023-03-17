package res

import (
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	//go:embed *
	Files embed.FS
)

const (
	APP_NAME    = "ServerBoxMonitor"
	APP_VERSION = "0.0.1"
)

var (
	ServerBoxShellFileName = "monitor.sh"
	ServerBoxDirPath       = filepath.Join(os.Getenv("HOME"), ".config", "server_box")
	ServerBoxShellPath     = filepath.Join(ServerBoxDirPath, ServerBoxShellFileName)

	AppConfigFileName = "config.json"
	AppConfigPath     = filepath.Join(ServerBoxDirPath, AppConfigFileName)
)

const (
	DefaultInterval = time.Minute
)

func init() {
	if !utils.Exist(ServerBoxDirPath) {
		err := os.MkdirAll(ServerBoxDirPath, 0755)
		if err != nil {
			utils.Error("[INIT] Create dir error: %v", err)
			panic(err)
		}
	}
}
