package model

import "github.com/lollipopkit/gommon/sys"

type WebConfig struct {
	Addr string `json:"addr"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func (wc *WebConfig) HaveTLS() bool {
	if wc == nil {
		return false
	}
	if len(wc.Cert) == 0 || len(wc.Key) == 0 {
		return false
	}
	if !sys.Exist(wc.Cert) || !sys.Exist(wc.Key) {
		return false
	}
	return true
}
