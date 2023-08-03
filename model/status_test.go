package model_test

import (
	_ "embed"
	"testing"

	"github.com/lollipopkit/server_box_monitor/model"
)

var (
	//go:embed test/disk
	_disk string
)

func TestParseDisk(t *testing.T) {
	err := model.ParseDiskStatus(_disk)
	if err != nil {
		t.Error(err)
	}
	t.Log(model.Status.Disk)
}
