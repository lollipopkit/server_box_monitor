package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lollipopkit/gommon/util"
)

type Push struct {
	Type      PushType        `json:"type"`
	Iface     json.RawMessage `json:"iface"`
	BodyRegex string          `json:"body_regex"`
	Code      int             `json:"code"`
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
	}
	return nil, errors.New("unknown push type")
}

func (p *Push) Push(args []*PushPair) error {
	iface, err := p.GetIface()
	if err != nil {
		return err
	}
	resp, code, err := iface.push(args)
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
func (p *Push) Id() string {
	iface, err := p.GetIface()
	if err != nil {
		return "UnknownPushIface"
	}
	switch iface.(type) {
	case PushIfaceIOS:
		return iface.(PushIfaceIOS).Name
	case PushIfaceWebhook:
		return iface.(PushIfaceWebhook).Name
	case PushIfaceServerChan:
		return iface.(PushIfaceServerChan).Name
	default:
		return fmt.Sprintf("UnknownPushId%v", iface)
	}
}

// {{key}} {{value}}
type PushFormat string
type PushPair struct {
	Key   string
	Value string
}

func (pf PushFormat) Format(args []*PushPair) string {
	ss := []string{}
	for _, arg := range args {
		s := string(pf)
		s = strings.ReplaceAll(s, "{{key}}", arg.Key)
		s = strings.ReplaceAll(s, "{{value}}", arg.Value)
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

type PushType string

const (
	PushTypeIOS        PushType = "ios"
	PushTypeWebhook             = "webhook"
	PushTypeServerChan          = "server_chan"
)

type PushIface interface {
	push([]*PushPair) ([]byte, int, error)
}

type PushIfaceIOS struct {
	Name    string     `json:"name"`
	Token   string     `json:"token"`
	Title   PushFormat `json:"title"`
	Content PushFormat `json:"content"`
}

func (p PushIfaceIOS) push(args []*PushPair) ([]byte, int, error) {
	title := p.Title.Format(args)
	content := p.Content.Format(args)
	body := map[string]string{
		"token":   p.Token,
		"title":   title,
		"content": content,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}
	return util.HttpDo(
		"POST",
		"https://push.lolli.tech/v1/ios",
		string(bodyBytes),
		map[string]string{
			"Content-Type": "application/json",
			"AppID":        "com.lollipopkit.toolbox",
		},
	)
}

type PushIfaceWebhook struct {
	Name    string            `json:"name"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Body    json.RawMessage   `json:"body"`
}

func (p PushIfaceWebhook) push(args []*PushPair) ([]byte, int, error) {
	body := PushFormat(p.Body).Format(args)
	switch p.Method {
	case "GET":
		return util.HttpDo("GET", p.Url, body, p.Headers)
	case "POST":
		return util.HttpDo("POST", p.Url, body, p.Headers)
	}
	return nil, 0, fmt.Errorf("unknown method: %s", p.Method)
}

type PushIfaceServerChan struct {
	Name  string     `json:"name"`
	SCKEY string     `json:"sckey"`
	Title PushFormat `json:"title"`
	Desp  PushFormat `json:"desp"`
}

func (p PushIfaceServerChan) push(args []*PushPair) ([]byte, int, error) {
	title := p.Title.Format(args)
	desp := p.Desp.Format(args)
	url := fmt.Sprintf("https://sctapi.ftqq.com/%s.send?title=%s&desp=%s", p.SCKEY, title, desp)
	return util.HttpDo(
		"GET",
		url,
		"",
		nil,
	)
}
