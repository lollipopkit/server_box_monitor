package model

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidRule   = errors.New("invalid rule")
	ErrInvalidMonitorType   = errors.New("invalid monitor type")
	ErrInvalidThresholdType = errors.New("invalid threshold type")
	ErrCompareFailed = errors.New("compare failed")

	ErrThresholdNotSupportPercent = errors.New("not support threshold type: percent, should use amount type, eg: >=80m")
)

type Rule struct {
	// eg: "cpu"
	MonitorType `json:"type"`
	// eg: ">=80.5%"
	Threshold string `json:"threshold"`
	// eg: "eth0" "cpu0" "sda1"
	Matcher string `json:"matcher"`
}

func (r *Rule) ShouldNotify(s *Status) (bool, error) {
	t, err := ParseToThreshold(r.Threshold)
	if err != nil {
		return false, errors.Join(ErrInvalidRule, err)
	}
	switch r.MonitorType {
	case MonitorTypeCPU:
		return r.shouldNotifyCPU(s.CPU, t)
	case MonitorTypeMemory:
		return r.shouldNotifyMemory(s.Mem, t)
	case MonitorTypeSwap:
		return r.shouldNotifySwap(s.Swap, t)
	case MonitorTypeDisk:
		return r.shouldNotifyDisk(s.Disk, t)
	case MonitorTypeNetwork:
		return r.shouldNotifyNetwork(s.Network, t)
	default:
		return false, errors.Join(ErrInvalidRule, ErrInvalidMonitorType)
	}
}

func (r *Rule) shouldNotifyCPU(ss []CPUStatus, t *Threshold) (bool, error) {
	idx, err := strconv.ParseInt(r.Matcher, 10, 64)
	if err != nil {
		return false, errors.Join(ErrInvalidRule, err)
	}
	if idx < 0 || int(idx) >= len(ss) {
		return false, errors.Join(ErrInvalidRule, fmt.Errorf("cpu index out of range: %d", idx))
	}
	s := ss[idx]
	switch t.ThresholdType {
	case ThresholdTypeSize:
		return false, errors.Join(ErrInvalidRule, ErrThresholdNotSupportPercent)
	case ThresholdTypePercent:
		return t.True(s.UsedPercent)
	default:
		return false, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifyMemory(s MemStatus, t *Threshold) (bool, error) {
	var size Size
	var percent float64
	switch r.Matcher {
	case "used":
		size = s.Used
		percent = float64(s.Used) / float64(s.Total)
	case "free":
		size = s.Free
		percent = float64(s.Free) / float64(s.Total)
	case "cached":
		size = s.Cached
		percent = float64(s.Cached) / float64(s.Total)
	default:
		return false, errors.Join(ErrInvalidRule, fmt.Errorf("invalid matcher: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		return t.True(size)
	case ThresholdTypePercent:
		return t.True(percent)
	default:
		return false, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifySwap(s SwapStatus, t *Threshold) (bool, error) {
	var size Size
	var percent float64
	switch r.Matcher {
	case "used":
		size = s.Used
		percent = float64(s.Used) / float64(s.Total)
	case "free":
		size = s.Free
		percent = float64(s.Free) / float64(s.Total)
	default:
		return false, errors.Join(ErrInvalidRule, fmt.Errorf("invalid matcher: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		return t.True(size)
	case ThresholdTypePercent:
		return t.True(percent)
	default:
		return false, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifyDisk(s []DiskStatus, t *Threshold) (bool, error) {
	var disk DiskStatus
	var have bool
	for _, d := range s {
		if d.Device == r.Matcher {
			disk = d
			have = true
			break
		}
	}
	if !have {
		return false, errors.Join(ErrInvalidRule, fmt.Errorf("disk not found: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		return t.True(disk.UsedAmount)
	case ThresholdTypePercent:
		return t.True(disk.UsedPercent)
	default:
		return false, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifyNetwork(s []NetworkStatus, t *Threshold) (bool, error) {
	var net NetworkStatus
	var have bool
	for _, n := range s {
		if n.Interface == r.Matcher {
			net = n
			have = true
			break
		}
	}
	if !have {
		return false, errors.Join(ErrInvalidRule, fmt.Errorf("network interface not found: %s", r.Matcher))
	}
	switch t.ThresholdType {
	case ThresholdTypeSize:
		return t.True(net.TransmitSpeed)
	case ThresholdTypePercent:
		return false, errors.Join(ErrInvalidRule, ErrThresholdNotSupportPercent)
	default:
		return false, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}

type MonitorType string

const (
	MonitorTypeCPU     MonitorType = "cpu"
	MonitorTypeMemory              = "mem"
	MonitorTypeSwap                = "swap"
	MonitorTypeDisk                = "disk"
	MonitorTypeNetwork             = "network"
)
