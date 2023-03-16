package res

import (
	"embed"
)

var (
	//go:embed *
	Files embed.FS
)

const (
	APP_NAME = "ServerBoxMonitor"
	APP_VERSION = "0.0.1"
)

var (
	ServerBoxShellFileName = "monitor.sh"
	ServerBoxDirPath = ".config/server_box/"
	// ServerBoxDirPath = os.Getenv("HOME") + ".config/server_box/"
	ServerBoxShellPath = ServerBoxDirPath + ServerBoxShellFileName

	AppConfigPath = ServerBoxDirPath + "config.json"
)
