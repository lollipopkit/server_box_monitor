package model

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/lollipopkit/gommon/http"
	"github.com/lollipopkit/server_box_monitor/res"
)

type PushType string

const (
	PushTypeIOS        PushType = "ios"
	PushTypeWebhook             = "webhook"
	PushTypeServerChan          = "server_chan"
	PushTypeBark                = "bark"
)

type Push struct {
	Type  PushType        `json:"type"`
	Name  string          `json:"name"`
	Iface json.RawMessage `json:"iface"`
}

func (p *Push) GetIface() (PushIface, error) {
	switch p.Type {
	case PushTypeIOS:
		var iface PushIfaceIOS
		err := json.Unmarshal(p.Iface, &iface)
		if err != nil {
			return nil, err
		}
		return iface, nil
	case PushTypeWebhook:
		var iface PushIfaceWebhook
		err := json.Unmarshal(p.Iface, &iface)
		if err != nil {
			return nil, err
		}
		return iface, nil
	case PushTypeServerChan:
		var iface PushIfaceServerChan
		err := json.Unmarshal(p.Iface, &iface)
		if err != nil {
			return nil, err
		}
		return iface, nil
	case PushTypeBark:
		var iface PushIfaceBark
		err := json.Unmarshal(p.Iface, &iface)
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("unknown push type: %s", p.Type)
}

func (p *Push) Push(args []*PushPair) error {
	iface, err := p.GetIface()
	if err != nil {
		return err
	}
	return iface.push(args)
}

type PushFormat string
type PushPair struct {
	key   string
	value string
	time  string
}

func NewPushPair(key, value string) *PushPair {
	return &PushPair{
		key:   key,
		value: value,
		time:  time.Now().Format("15:04:05"),
	}
}

func (pf PushFormat) Format(args []*PushPair) string {
	ss := []string{}
	for _, arg := range args {
		if len(arg.time) == 0 {
			arg.time = time.Now().Format("15:04:05")
		}
		kv := fmt.Sprintf(`%s\n%s: %s`, arg.time, arg.key, arg.value)
		ss = append(ss, kv)
	}
	msgReplaced := strings.Replace(
		string(pf),
		res.PushFormatMsgLocator,
		strings.Join(ss, `\n`),
		1,
	)
	nameReplaced := strings.Replace(
		msgReplaced,
		res.PushFormatNameLocator,
		Config.Name,
		1,
	)
	return nameReplaced
}

type PushIface interface {
	push([]*PushPair) error
}

type PushIfaceIOS struct {
	Token     string     `json:"token"`
	Title     PushFormat `json:"title"`
	Content   PushFormat `json:"content"`
	BodyRegex string     `json:"body_regex"`
	Code      int        `json:"code"`
}

func (p PushIfaceIOS) push(args []*PushPair) error {
	content := p.Content.Format(args)
	title := p.Title.Format(args)
	body := map[string]string{
		"token":   p.Token,
		"title":   title,
		"content": content,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, code, err := http.Do(
		"POST",
		"https://push.lolli.tech/v1/ios",
		bodyBytes,
		map[string]string{
			"AppID":        "com.lollipopkit.toolbox",
			"Content-Type": "application/json",
			"Env":          "prod",
		},
	)

	if p.Code != 0 && code != p.Code {
		return fmt.Errorf("code: %d, resp: %s", code, string(resp))
	}
	if p.BodyRegex != "" {
		reg, err := regexp.Compile(p.BodyRegex)
		if err != nil {
			return fmt.Errorf("compile regex failed: %s", err.Error())
		}
		if !reg.Match(resp) {
			return fmt.Errorf("resp: %s", string(resp))
		}
	}
	return nil
}

type PushIfaceWebhook struct {
	Url       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Method    string            `json:"method"`
	Body      json.RawMessage   `json:"body"`
	BodyRegex string            `json:"body_regex"`
	Code      int               `json:"code"`
}

func (p PushIfaceWebhook) push(args []*PushPair) error {
	body := PushFormat(p.Body).Format(args)
	switch p.Method {
	case "GET", "POST":
		resp, code, err := http.Do(p.Method, p.Url, body, p.Headers)
		if err != nil {
			return err
		}
		if p.Code != 0 && code != p.Code {
			return fmt.Errorf("code: %d, resp: %s", code, string(resp))
		}
		if p.BodyRegex != "" {
			reg, err := regexp.Compile(p.BodyRegex)
			if err != nil {
				return fmt.Errorf("compile regex failed: %s", err.Error())
			}
			if !reg.Match(resp) {
				return fmt.Errorf("resp: %s", string(resp))
			}
		}
		return nil
	}
	return fmt.Errorf("unknown method: %s", p.Method)
}

type PushIfaceServerChan struct {
	SCKey     string     `json:"sckey"`
	Title     PushFormat `json:"title"`
	Desp      PushFormat `json:"desp"`
	BodyRegex string     `json:"body_regex"`
	Code      int        `json:"code"`
}

func (p PushIfaceServerChan) push(args []*PushPair) error {
	desp := p.Desp.Format(args)
	title := p.Title.Format(args)
	url := fmt.Sprintf(
		"https://sctapi.ftqq.com/%s.send?title=%s&desp=%s",
		p.SCKey,
		title,
		desp,
	)
	resp, code, err := http.Do("GET", url, nil, nil)
	if err != nil {
		return err
	}
	if p.Code != 0 && code != p.Code {
		return fmt.Errorf("code: %d, resp: %s", code, string(resp))
	}
	if p.BodyRegex != "" {
		reg, err := regexp.Compile(p.BodyRegex)
		if err != nil {
			return fmt.Errorf("compile regex failed: %s", err.Error())
		}
		if !reg.Match(resp) {
			return fmt.Errorf("resp: %s", string(resp))
		}
	}
	return nil
}

type barkLevel string

const (
	barkLevelActive    barkLevel = "active"
	barkLevelSensitive           = "timeSensitive"
	barkLevelPassive             = "passive"
)

type PushIfaceBark struct {
	Server    string    `json:"server"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Level     barkLevel `json:"level"`
	BodyRegex string    `json:"body_regex"`
	Code      int       `json:"code"`
}

func (p PushIfaceBark) push(args []*PushPair) error {
	body := p.Body
	for _, arg := range args {
		body = strings.Replace(body, arg.key, arg.value, 1)
	}
	if len(p.Server) == 0 {
		p.Server = "https://api.day.app"
	}
	url := path.Join(
		p.Server,
		p.Title,
		body,
	)
	if len(p.Level) != 0 {
		url += fmt.Sprintf("?level=%s", p.Level)
	}
	resp, code, err := http.Do("GET", url, nil, nil)
	if err != nil {
		return err
	}
	if p.Code != 0 && code != p.Code {
		return fmt.Errorf("code: %d, resp: %s", code, string(resp))
	}
	if p.BodyRegex != "" {
		reg, err := regexp.Compile(p.BodyRegex)
		if err != nil {
			return fmt.Errorf("compile regex failed: %s", err.Error())
		}
		if !reg.Match(resp) {
			return fmt.Errorf("resp: %s", string(resp))
		}
	}
	return nil
}
