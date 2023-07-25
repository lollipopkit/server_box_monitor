package res

import (
	"embed"
	"os"
	"path/filepath"
	"time"

	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/rate"
	"github.com/lollipopkit/gommon/sys"
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

	DefaultRateLimiter = rate.NewLimiter[string](time.Second*10, 1)
)

const (
	ConfVersion = 2

	DefaultInterval    = time.Second * 7
	DefaultIntervalStr = "7s"
	DefaultRateStr     = "1/1m"
	DefaultSeverName   = "Server 1"
	MaxInterval        = time.Second * 10

	PushFormatMsgLocator  = "{{msg}}"
	PushFormatNameLocator = "{{name}}"
)

func init() {
	if !sys.Exist(ServerBoxDirPath) {
		err := os.MkdirAll(ServerBoxDirPath, 0755)
		if err != nil {
			log.Err("[INIT] Create dir error: %v", err)
			panic(err)
		}
	}
}
