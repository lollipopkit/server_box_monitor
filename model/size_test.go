package model_test

import (
	"testing"

	"github.com/lollipopkit/server_box_monitor/model"
)

func TestParseToSize(t *testing.T) {
	_parseSize("1m", model.Size(1024*1024), t)
	_parseSize("1M", model.Size(1024*1024), t)
	_parseSize("3k", model.Size(3*1024), t)
	_parseSize("7b", model.Size(7), t)
}

func _parseSize(s string, expect model.Size, t *testing.T) {
	size, err := model.ParseToSize(s)
	if err != nil {
		t.Error(err)
	}
	if size != expect {
		t.Errorf("expect %s, got %s", expect, size)
	}
}
