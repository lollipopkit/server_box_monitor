package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	ErrInvalidRule          = errors.New("invalid rule")
	ErrInvalidMonitorType   = errors.New("invalid monitor type")
	ErrInvalidThresholdType = errors.New("invalid threshold type")
	ErrCompareFailed        = errors.New("compare failed")
)

type Rule struct {
	// eg: "cpu"
	MonitorType `json:"type"`
	// eg: ">=80.5%" "<100m" "=10m/s"
	// Threshold which match speed should use per second
	// such as "10m/s"
	// "10m/m" is not allowed
	Threshold string `json:"threshold"`
	// eg: "eth0-in" "cpu0" "sda1" "free"
	// "cpu0" -> all CPUs
	// MonitorType = "mem" && Matcher = "free" -> free of memory
	// MonitorType = "net" && Matcher = "eth0-in" -> in speed of eth0
	// MonitorType = "net" && Matcher = "eth0-out-in" -> out + in speed of eth0
	// MonitorType = "disk" && Matcher = "sda1" -> used percent of sda1
	Matcher string `json:"matcher"`
}

func (r *Rule) Id() string {
	return fmt.Sprintf("%s-%s-%s", r.MonitorType, r.Threshold, r.Matcher)
}
func (r *Rule) ShouldNotify(s *Status) (bool, *PushFormatArgs, error) {
	t, err := ParseToThreshold(r.Threshold)
	if err != nil {
		return false, nil, errors.Join(ErrInvalidRule, err)
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
	case MonitorTypeTemperature:
		return r.shouldNotifyTemperature(s.Temperature, t)
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidMonitorType)
	}
}

func (r *Rule) shouldNotifyCPU(ss []CPUStatus, t *Threshold) (bool, *PushFormatArgs, error) {
	if len(ss) == 0 {
		// utils.Warn("cpu is not valid, skip this rule")
		return false, nil, nil
	}
	idx, err := strconv.ParseInt(strings.Replace(r.Matcher, "cpu", "", 1), 10, 64)
	if err != nil {
		return false, nil, errors.Join(ErrInvalidRule, err)
	}
	if idx < 0 || int(idx) >= len(ss) {
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("cpu index out of range: %d", idx))
	}
	s := ss[idx]
	switch t.ThresholdType {
	case ThresholdTypePercent:
		ok, err := t.True(s.UsedPercent)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   fmt.Sprintf("cpu%d", idx),
			Value: fmt.Sprintf("%.2f%%", s.UsedPercent),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifyMemory(s *MemStatus, t *Threshold) (bool, *PushFormatArgs, error) {
	if s == nil {
		// utils.Warn("memory is not valid, skip this rule")
		return false, nil, nil
	}
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
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid matcher: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		ok, err := t.True(size)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher + "of Memory",
			Value: size.String(),
		}, nil
	case ThresholdTypePercent:
		ok, err := t.True(percent)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher + "of Memory",
			Value: fmt.Sprintf("%.2f%%", percent*100),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifySwap(s *SwapStatus, t *Threshold) (bool, *PushFormatArgs, error) {
	if s == nil {
		// utils.Warn("swap is not valid, skip this rule")
		return false, nil, nil
	}
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
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid matcher: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		ok, err := t.True(size)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher + "of Swap",
			Value: size.String(),
		}, nil
	case ThresholdTypePercent:
		ok, err := t.True(percent)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher + "of Swap",
			Value: fmt.Sprintf("%.2f%%", percent*100),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifyDisk(s []DiskStatus, t *Threshold) (bool, *PushFormatArgs, error) {
	if len(s) == 0 {
		// utils.Warn("disk is not valid, skip this rule")
		return false, nil, nil
	}
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
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("disk not found: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		ok, err := t.True(disk.UsedAmount)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher,
			Value: disk.UsedAmount.String(),
		}, nil
	case ThresholdTypePercent:
		ok, err := t.True(disk.UsedPercent)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher,
			Value: fmt.Sprintf("%.2f%%", disk.UsedPercent),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}
func (r *Rule) shouldNotifyNetwork(s []NetworkStatus, t *Threshold) (bool, *PushFormatArgs, error) {
	if len(s) == 0 {
		// utils.Warn("network is not valid, skip this rule")
		return false, nil, nil
	}

	var net NetworkStatus
	var have bool
	for _, n := range s {
		if strings.Contains(r.Matcher, n.Interface)  {
			net = n
			have = true
			break
		}
	}
	if !have {
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("network interface not found: %s", r.Matcher))
	}

	// 判断是否计算出/入流量
	in := strings.Contains(r.Matcher, "-in")
	out := strings.Contains(r.Matcher, "-out")
	if !in && !out {
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid matcher: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		speed := Size(0)
		if in {
			s, err := net.ReceiveSpeed()
			if err != nil {
				utils.Warn("[NETWORK] get receive speed failed: %s", err)
			}
			speed += s
		}
		if out {
			s, err := net.TransmitSpeed()
			if err != nil {
				utils.Warn("[NETWORK] get transmit speed failed: %s", err)
			}
			speed += s
		}
		ok, err := t.True(speed)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher,
			Value: speed.String(),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}

func (r *Rule) shouldNotifyTemperature(s []TemperatureStatus, t *Threshold) (bool, *PushFormatArgs, error) {
	if len(s) == 0 {
		// utils.Warn("temperature is not valid, skip this rule")
		return false, nil, nil
	}

	var temp TemperatureStatus
	var have bool
	for _, t := range s {
		if strings.Contains(t.Name, r.Matcher) {
			temp = t
			have = true
			break
		}
	}
	if !have {
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("temperature not found: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		ok, err := t.True(temp.Value)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		return ok, &PushFormatArgs{
			Key:   r.Matcher,
			Value: fmt.Sprintf("%.2f°C", temp.Value),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, ErrInvalidThresholdType)
	}
}

type MonitorType string

const (
	MonitorTypeCPU     MonitorType = "cpu"
	MonitorTypeMemory              = "mem"
	MonitorTypeSwap                = "swap"
	MonitorTypeDisk                = "disk"
	MonitorTypeNetwork             = "net"
	MonitorTypeTemperature         = "temp"
)
