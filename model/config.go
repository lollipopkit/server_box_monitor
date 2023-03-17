package model

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	Config = &AppConfig{}
)

type AppConfig struct {
	Version  int `json:"version"`
	Interval `json:"interval"`
	Rules    []Rule `json:"rules"`
	Pushes   []Push `json:"pushes"`
}

// Such as "300ms", "-1.5h" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
type Interval struct {
	// interval for running the script to get status
	// Values greater than 10 or less than 3 will be ignored
	Run string `json:"run"`
	// interval for pushing the status
	Push string `json:"push"`
}

var (
	DefaultappConfig = &AppConfig{
		Version:  1,
		Interval: Interval{Run: "5s", Push: "5m"},
		Rules: []Rule{
			{
				MonitorType: MonitorTypeCPU,
				Threshold:   ">=77%",
				Matcher:     "0",
			},
			{
				MonitorType: MonitorTypeNetwork,
				Threshold:   ">=7.7m/s",
				Matcher:     "eth0",
			},
		},
		Pushes: []Push{
			{
				Type: PushTypeWebhook,
				Iface: []byte(`{
					"name": "QQ Group",
					"url": "http://localhost:5700",
					"headers": {
						"Content-Type": "application/json"
						"Auhtorization": "Bearer YOUR_SECRET"
					},
					"method": "POST",
					"body": {
						"action": "send_group_msg",
						"params": {
							"group_id": 123456789,
							"message": "ServerBox Notification: {{key}}: {{value}}"
						}
					}
				}`),
				SuccessBodyRegex: ".*",
				SuccessCode:      200,
			},
			{
				Type: PushTypeIOS,
				Iface: []byte(`{
					"name": "My iPhone",
					"token": "YOUR_TOKEN",
					"title": "Server Notification",
					"content": "{{key}}: {{value}}"
				}`),
				SuccessBodyRegex: ".*",
				SuccessCode:      200,
			},
		},
	}
)

func ReadAppConfig() error {
	if !utils.Exist(res.AppConfigPath) {
		configBytes, err := json.MarshalIndent(DefaultappConfig, "", "\t")
		if err != nil {
			utils.Error("[CONFIG] marshal default app config failed: %v", err)
			return err
		}
		err = os.WriteFile(res.AppConfigPath, configBytes, 0644)
		if err != nil {
			utils.Error("[CONFIG] write default app config failed: %v", err)
			return err
		}
		Config = DefaultappConfig
		return nil
	}

	configBytes, err := os.ReadFile(res.AppConfigPath)
	if err != nil {
		utils.Error("[CONFIG] read app config failed: %v", err)
		return err
	}
	err = json.Unmarshal(configBytes, Config)
	if err != nil {
		utils.Error("[CONFIG] unmarshal app config failed: %v", err)
	} else if Config.Version < DefaultappConfig.Version {
		utils.Warn("[CONFIG] app config version is too old, please update it")
	}
	return err
}

func (ac *AppConfig) GetRunInterval() time.Duration {
	d, err := time.ParseDuration(ac.Interval.Run)
	if err == nil {
		return d
	}
	utils.Warn("[CONFIG] parse interval failed: %v, use default interval: 5s", err)
	return time.Second * 5
}
func (ac *AppConfig) GetPushInterval() time.Duration {
	d, err := time.ParseDuration(ac.Interval.Push)
	if err == nil {
		return d
	}
	utils.Warn("[CONFIG] parse interval failed: %v, use default interval: 5m", err)
	return time.Minute * 5
}
