package model

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/lollipopkit/server_box_monitor/utils"
)

type Push struct {
	Type          PushType         `json:"type"`
	Iface         json.RawMessage `json:"iface"`
	TitleFormat   PushFormat       `json:"title"`
	ContentFormat PushFormat       `json:"content"`
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

func (p *Push) Push(args []*PushFormatArgs) error {
	title := p.TitleFormat.String(args)
	content := p.ContentFormat.String(args)
	iface, err := p.GetIface()
	if err != nil {
		return err
	}
	return iface.push(title, content)
}
func (p *Push) Id() string {
	// switch p.PushType {
	// case PushTypeIOS:
	// 	return "iOS-" + p.PushIface.(*PushIOS).Token[:7]
	// case PushTypeWebhook:
	// 	return "Webhook-" + p.PushIface.(*PushWebhook).Url
	// default:
	// 	return "UnknownPushId"
	// }
	iface, err := p.GetIface()
	if err != nil {
		return "UnknownPushIface"
	}
	switch iface.(type) {
	case PushIOS:
		return "iOS-" + iface.(PushIOS).Token[:7]
	case PushWebhook:
		return "Webhook-" + iface.(PushWebhook).Url
	default:
		return "UnknownPushId"
	}
}

// {{key}} {{value}}
type PushFormat string
type PushFormatArgs struct {
	Key   string
	Value string
}

func (pf PushFormat) String(args []*PushFormatArgs) string {
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
	push(title, content string) error
}

type PushIOS struct {
	Token string `json:"token"`
}

func (p PushIOS) push(title, content string) error {
	return nil
}

type PushWebhook struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
}

func (p PushWebhook) push(title, content string) error {
	body := strings.Join([]string{title, content}, "\n")
	switch p.Method {
	case "GET":
		resp, err := utils.HttpDo("GET", p.Url, body, p.Headers)
		if err != nil {
			utils.Warn("[PUSH] webhook GET failed: '%v', resp: '%s'", err, resp)
		}
	case "POST":
		resp, err := utils.HttpDo("POST", p.Url, body, p.Headers)
		if err != nil {
			utils.Warn("[PUSH] webhook POST failed: '%v', resp: '%s'", err, resp)
		}
	}
	return nil
}
