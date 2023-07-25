package model

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/rate"
	"github.com/lollipopkit/gommon/sys"
	"github.com/lollipopkit/server_box_monitor/res"
)

var (
	Config        = new(AppConfig)
	CheckInterval time.Duration
	RateLimiter   *rate.RateLimiter[string]
)

type AppConfig struct {
	Version int `json:"version"`
	// Such as "7s".
	// Valid time units are "s".
	// Values bigger than 10 seconds are not allowed.
	Interval string `json:"interval"`
	Rate     string `json:"rate"`
	Name     string `json:"name"`
	Rules    []Rule `json:"rules"`
	Pushes   []Push `json:"pushes"`
}

func InitConfig() error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err := enc.Encode(DefaultAppConfig)
	if err != nil {
		log.Err("[CONFIG] marshal default app config failed: %v", err)
		return err
	}
	err = os.WriteFile(res.AppConfigPath, buf.Bytes(), 0644)
	if err != nil {
		log.Err("[CONFIG] write default app config failed: %v", err)
		return err
	}
	Config = DefaultAppConfig
	return nil
}

func ReadAppConfig() error {
	defer initInterval()
	defer initRateLimiter()
	if !sys.Exist(res.AppConfigPath) {
		return InitConfig()
	}

	configBytes, err := os.ReadFile(res.AppConfigPath)
	if err != nil {
		log.Err("[CONFIG] read app config failed: %v", err)
		return err
	}
	err = json.Unmarshal(configBytes, Config)
	if err != nil {
		log.Err("[CONFIG] unmarshal app config failed: %v", err)
	} else if Config.Version < DefaultAppConfig.Version {
		log.Warn("[CONFIG] app config version is too old, new config will be generated")
		// Backup old config
		err = os.WriteFile(res.AppConfigPath+".bak", configBytes, 0644)
		if err != nil {
			log.Err("[CONFIG] backup old config failed: %v", err)
			return err
		}
		// Generate new config
		configBytes, err := json.MarshalIndent(DefaultAppConfig, "", "\t")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(res.AppConfigPath, configBytes, 0644)
		if err != nil {
			panic(err)
		}
		log.Info("[CONFIG] new config generated, edit it and restart the program")
		os.Exit(0)
	}
	return err
}

func initInterval() {
	d, err := time.ParseDuration(Config.Interval)
	if err == nil {
		if d > res.MaxInterval || d < time.Second {
			log.Warn("[CONFIG] use default interval")
			CheckInterval = res.DefaultInterval
			return
		}
		CheckInterval = d
		return
	}
	log.Warn("[CONFIG] parse interval failed: %v", err)
	CheckInterval = res.DefaultInterval
}

func initRateLimiter() {
	splited := strings.Split(Config.Rate, "/")
	if len(splited) != 2 {
		log.Warn("[CONFIG] parse rate failed")
		RateLimiter = res.DefaultRateLimiter
		return
	}
	times, err := strconv.Atoi(splited[0])
	if err != nil {
		log.Warn("[CONFIG] parse rate failed: %v", err)
		RateLimiter = res.DefaultRateLimiter
		return
	}
	duration, err := time.ParseDuration(splited[1])
	if err != nil {
		log.Warn("[CONFIG] parse rate failed: %v", err)
		RateLimiter = res.DefaultRateLimiter
		return
	}
	RateLimiter = rate.NewLimiter[string](duration, times)
}

var (
	defaultWebhookBody = map[string]interface{}{
		"action": "send_group_msg",
		"params": map[string]interface{}{
			"group_id": 123456789,
			"message": res.PushFormatNameLocator +
				"\n" +
				res.PushFormatMsgLocator,
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

	DefaultAppConfig = &AppConfig{
		Version:  res.ConfVersion,
		Interval: res.DefaultIntervalStr,
		Rate:     res.DefaultRateStr,
		Name:     res.DefaultSeverName,
		Rules: []Rule{
			{
				MonitorType: MonitorTypeCPU,
				Threshold:   `>=77%`,
				Matcher:     "cpu",
			},
		},
		Pushes: []Push{
			{
				Type:  PushTypeWebhook,
				Name:  "QQ Group",
				Iface: defaultWebhookIfaceBytes,
			},
		},
	}
)
