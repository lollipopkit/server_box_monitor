package model

import (
	"fmt"
	"strings"

	"github.com/lollipopkit/server_box_monitor/utils"
)

type Push struct {
	PushType `json:"type"`
	PushIface `json:"iface"`
	TitleFormat PushFormat `json:"title"`
	ContentFormat PushFormat `json:"content"`
}
func (p *Push) Push(args []PushFormatArgs) error {
	title := p.TitleFormat.String(args)
	content := p.ContentFormat.String(args)
	return p.PushIface.push(title, content)
}

// 支持的格式化参数
// {{key}}: cpu, mem, swap, disk, network
// {{value}}: 80.5%
type PushFormat string
type PushFormatArgs struct {
	Key string
	Value string
}
func (pf PushFormat) String(args []PushFormatArgs) string {
	s := string(pf)
	for _, arg := range args {
		key := fmt.Sprintf("{{%s}}", arg.Key)
		pf = PushFormat(strings.ReplaceAll(s, key, arg.Value))
	}
	return s
}

type PushType string
const (
	PushTypeIOS PushType = "ios"
	PushTypeWebhook = "webhook"
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
	Url string `json:"url"`
	Headers map[string]string `json:"headers"`
	Method string `json:"method"`
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