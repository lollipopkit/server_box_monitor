package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/lollipopkit/server_box_monitor/utils"
)

type Push struct {
	Type             PushType        `json:"type"`
	Iface            json.RawMessage `json:"iface"`
	SuccessBodyRegex string          `json:"success_body_regex"`
	SuccessCode      int             `json:"success_code"`
}

func (p *Push) GetIface() (PushIface, error) {
	switch p.Type {
	case PushTypeIOS:
		var iface PushIOS
		err := json.Unmarshal(p.Iface, &iface)
		if err != nil {
			return nil, err
		}
		return iface, nil
	case PushTypeWebhook:
		var iface PushWebhook
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
	if p.SuccessCode != 0 && code != p.SuccessCode {
		return fmt.Errorf("code: %d, resp: %s", code, string(resp))
	}
	if p.SuccessBodyRegex != "" {
		reg, err := regexp.Compile(p.SuccessBodyRegex)
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
	case PushIOS:
		return iface.(PushIOS).Name
	case PushWebhook:
		return iface.(PushWebhook).Name
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

func (pf PushFormat) String(args []*PushPair) string {
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
	PushTypeIOS     PushType = "ios"
	PushTypeWebhook          = "webhook"
)

type PushIface interface {
	push([]*PushPair) ([]byte, int, error)
}

type PushIOS struct {
	Name string `json:"name"`
	Token string `json:"token"`
	Title PushFormat `json:"title"`
	Content PushFormat `json:"content"`
}

func (p PushIOS) push(args []*PushPair) ([]byte, int, error) {
	title := p.Title.String(args)
	content := p.Content.String(args)
	func (a,b string){}(title, content)
	return nil, 0, errors.New("ios push now is not implemented")
}

type PushWebhook struct {
	Name string `json:"name"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Body json.RawMessage `json:"body"`
}

func (p PushWebhook) push(args []*PushPair) ([]byte, int, error) {
	body := PushFormat(p.Body).String(args)
	switch p.Method {
	case "GET":
		return utils.HttpDo("GET", p.Url, body, p.Headers)
	case "POST":
		return utils.HttpDo("POST", p.Url, body, p.Headers)
	}
	return nil, 0, fmt.Errorf("unknown method: %s", p.Method)
}
