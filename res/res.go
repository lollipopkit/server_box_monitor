package res

import (
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/util"
)

var (
	//go:embed *
	Files embed.FS
)

const (
	APP_NAME    = "ServerBoxMonitor"
	APP_VERSION = "0.0.4"
)

var (
	ServerBoxShellFileName = "monitor.sh"
	ServerBoxDirPath       = filepath.Join(os.Getenv("HOME"), ".config", "server_box")
	ServerBoxShellPath     = filepath.Join(ServerBoxDirPath, ServerBoxShellFileName)

	AppConfigFileName = "config.json"
	AppConfigPath     = filepath.Join(ServerBoxDirPath, AppConfigFileName)
)

const (
	DefaultInterval = time.Second * 7
)

func init() {
	if !util.Exist(ServerBoxDirPath) {
		err := os.MkdirAll(ServerBoxDirPath, 0755)
		if err != nil {
			log.Err("[INIT] Create dir error: %v", err)
			panic(err)
		}
	}
}
