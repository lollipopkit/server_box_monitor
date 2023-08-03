package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

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
		return iface, nil
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
}

func NewPushPair(key, value string) *PushPair {
	return &PushPair{
		key:   key,
		value: value,
	}
}

func (pf PushFormat) Format(args []*PushPair, raw bool) string {
	newline := `\n`
	if !raw {
		newline = "\n"
	}
	ss := []string{}
	for _, arg := range args {
		kv := fmt.Sprintf("%s: %s", arg.key, arg.value)
		ss = append(ss, kv)
	}
	msgReplaced := strings.Replace(
		string(pf),
		res.PushFormatMsgLocator,
		strings.Join(ss, newline),
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
	content := p.Content.Format(args, false)
	title := p.Title.Format(args, true)
	body := map[string]string{
		"token":   p.Token,
		"title":   title,
		"content": content,
	}
	resp, code, err := http.Do(
		"POST",
		"https://push.lolli.tech/v1/ios",
		body,
		map[string]string{
			"AppID":        "com.lollipopkit.toolbox",
			"Content-Type": "application/json",
			"Env":          "prod",
		},
	)

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

type PushIfaceWebhook struct {
	Url       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Method    string            `json:"method"`
	Body      json.RawMessage   `json:"body"`
	BodyRegex string            `json:"body_regex"`
	Code      int               `json:"code"`
}

func (p PushIfaceWebhook) push(args []*PushPair) error {
	body := PushFormat(p.Body).Format(args, true)
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
	desp := p.Desp.Format(args, true)
	title := p.Title.Format(args, true)
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
	Server    string     `json:"server"`
	Key       string     `json:"key"`
	Title     PushFormat `json:"title"`
	Body      PushFormat `json:"body"`
	Level     barkLevel  `json:"level"`
	BodyRegex string     `json:"body_regex"`
	Code      int        `json:"code"`
}

func (p PushIfaceBark) push(args []*PushPair) error {
	body := p.Body.Format(args, false)
	title := p.Title.Format(args, true)
	if len(p.Server) == 0 {
		p.Server = "https://api.day.app"
	}
	if strings.HasSuffix("/", p.Server) {
		p.Server = p.Server[:len(p.Server)-1]
	}
	titleEscape := url.QueryEscape(title)
	bodyEscape := url.QueryEscape(body)
	url_ := fmt.Sprintf(
		"%s/%s/%s/%s",
		p.Server, p.Key, titleEscape, bodyEscape,
	)
	resp, code, err := http.Do("GET", url_, nil, nil)
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
