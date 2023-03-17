package model

import (
	"strings"

	"github.com/lollipopkit/server_box_monitor/utils"
)

type Push struct {
	PushType      `json:"type"`
	PushIface     any        `json:"iface"`
	TitleFormat   PushFormat `json:"title"`
	ContentFormat PushFormat `json:"content"`
}

func (p *Push) Push(args []*PushFormatArgs) error {
	title := p.TitleFormat.String(args)
	content := p.ContentFormat.String(args)
	var iface PushIface
	switch p.PushType {
	case PushTypeIOS:
		iface = p.PushIface.(*PushIOS)
	case PushTypeWebhook:
		iface = p.PushIface.(*PushWebhook)
	}
	return iface.push(title, content)
}
func (p *Push) Id() string {
	switch p.PushType {
	case PushTypeIOS:
		return "iOS-" + p.PushIface.(*PushIOS).Token[:7]
	case PushTypeWebhook:
		return "Webhook-" + p.PushIface.(*PushWebhook).Url
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

func (p *PushIOS) push(title, content string) error {
	return nil
}

type PushWebhook struct {
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
}

func (p *PushWebhook) push(title, content string) error {
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
