package model

import (
	"errors"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

type Status struct {
	CPU     []CPUStatus
	Mem     MemStatus
	Swap    SwapStatus
	Disk    []DiskStatus
	Network []NetworkStatus
}

// All CPUs as one status
type CPUStatus struct {
	Core        int
	UsedPercent float64
}

type MemStatus struct {
	Total  Size
	Used   Size
	Free   Size
	Cached Size
}

type SwapStatus struct {
	Total Size
	Free  Size
	Used  Size
}

type DiskStatus struct {
	Device      string
	UsedPercent float64
	UsedAmount  Size
}

type NetworkStatus struct {
	Interface     string
	TransmitSpeed Size
	ReceiveSpeed  Size
}

var (
	ErrRunShellFailed = errors.New("run shell failed")
)

func GetStatus() (*Status, error) {
	stdout, stderr, err := utils.Execute("bash", res.ServerBoxShellPath)
	if err != nil {
		return nil, errors.Join(ErrRunShellFailed, errors.New(stderr))
	}
	return ParseStatus(stdout)
}

func ParseStatus(s string) (*Status, error) {
	return nil, ErrRunShellFailed
}
