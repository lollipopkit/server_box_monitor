package model

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/util"
	"github.com/lollipopkit/server_box_monitor/res"
)

var (
	Config = &AppConfig{}
)

type AppConfig struct {
	Version int `json:"version"`
	// Such as "7s".
	// Valid time units are "s".
	// Values bigger than 10 seconds are not allowed.
	Interval string `json:"interval"`
	Rules    []Rule `json:"rules"`
	Pushes   []Push `json:"pushes"`
}

func ReadAppConfig() error {
	if !util.Exist(res.AppConfigPath) {
		configBytes, err := json.MarshalIndent(DefaultappConfig, "", "\t")
		if err != nil {
			log.Err("[CONFIG] marshal default app config failed: %v", err)
			return err
		}
		err = os.WriteFile(res.AppConfigPath, configBytes, 0644)
		if err != nil {
			log.Err("[CONFIG] write default app config failed: %v", err)
			return err
		}
		Config = DefaultappConfig
		return nil
	}

	configBytes, err := os.ReadFile(res.AppConfigPath)
	if err != nil {
		log.Err("[CONFIG] read app config failed: %v", err)
		return err
	}
	err = json.Unmarshal(configBytes, Config)
	if err != nil {
		log.Err("[CONFIG] unmarshal app config failed: %v", err)
	} else if Config.Version < DefaultappConfig.Version {
		log.Warn("[CONFIG] app config version is too old, please update it")
	}
	return err
}

func GetInterval() time.Duration {
	ac := DefaultappConfig
	if Config != nil {
		ac = Config
	}
	d, err := time.ParseDuration(ac.Interval)
	if err == nil {
		if d > res.DefaultInterval {
			log.Warn("[CONFIG] interval is too long, use default interval")
			return res.DefaultInterval
		}
		return d
	}
	log.Warn("[CONFIG] parse interval failed: %v", err)
	return res.DefaultInterval
}

func GetIntervalInSeconds() float64 {
	return GetInterval().Seconds()
}

var (
	defaultWebhookBody = map[string]interface{}{
		"action": "send_group_msg",
		"params": map[string]interface{}{
			"group_id": 123456789,
			"message":  "ServerBox Notification\n{{key}}: {{value}}",
		},
	}
	defaultWekhookBodyBytes, _ = json.Marshal(defaultWebhookBody)
	defaultWebhookIface        = PushIfaceWebhook{
		Url: "http://localhost:5700",
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer YOUR_SECRET",
		},
		Method:    "POST",
		Body:      defaultWekhookBodyBytes,
		BodyRegex: ".*",
		Code:      200,
	}
	defaultWebhookIfaceBytes, _ = json.Marshal(defaultWebhookIface)

	defaultIOSIface = PushIfaceIOS{
		Token:     "YOUR_TOKEN",
		Title:     "Server Notification",
		Content:   "{{key}}: {{value}}",
		BodyRegex: ".*",
		Code:      200,
	}
	defaultIOSIfaceBytes, _ = json.Marshal(defaultIOSIface)

	defaultServerChanIface = PushIfaceServerChan{
		SCKey:     "YOUR_SCKEY",
		Title:     "Server Notification",
		Desp:      "{{key}}: {{value}}",
		BodyRegex: ".*",
		Code:      200,
	}
	defaultServerChanIfaceBytes, _ = json.Marshal(defaultServerChanIface)

	DefaultappConfig = &AppConfig{
		Version:  1,
		Interval: "30s",
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
				Type:  PushTypeWebhook,
				Name:  "QQ Group",
				Iface: defaultWebhookIfaceBytes,
			},
			{
				Type:  PushTypeIOS,
				Name:  "My iPhone",
				Iface: defaultIOSIfaceBytes,
			},
			{
				Type:  PushTypeServerChan,
				Name:  "ServerChan",
				Iface: defaultServerChanIfaceBytes,
			},
		},
	}
)
