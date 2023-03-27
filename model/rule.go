package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidRule = errors.New("invalid rule")
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
	// "cpu" -> all CPUs
	// MonitorType = "mem" && Matcher = "free" -> free of memory
	// MonitorType = "net" && Matcher = "eth0-in" -> in speed of eth0
	// MonitorType = "net" && Matcher = "eth0" -> out + in speed of eth0
	// MonitorType = "disk" && Matcher = "/dev/sda1" -> used percent of sda1
	// MonitorType = "disk" && Matcher = "/" -> used percent of mounted path "/"
	// MonitorType = "temp" && Matcher = "x86_pkg" -> temperature of x86_pkg
	Matcher string `json:"matcher"`
}

func (r *Rule) Id() string {
	return fmt.Sprintf("[%s %s %s]", r.MonitorType, r.Threshold, r.Matcher)
}
func (r *Rule) ShouldNotify(s *ServerStatus) (bool, *PushPair, error) {
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
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid monitor type: %s", r.MonitorType))
	}
}

func (r *Rule) shouldNotifyCPU(ss []OneCpuStatus, t *Threshold) (bool, *PushPair, error) {
	if len(ss) == 0 {
		// utils.Warn("cpu is not valid, skip this rule")
		return false, nil, nil
	}
	// 默认获取所有cpu
	// cpu -> idx = 0 （默认）
	// cpu0 -> idx = 1
	// idx = CPU序号 + 1
	var idx int64 = 0
	if r.Matcher != "" && r.Matcher != "cpu" {
		idx_, err := strconv.ParseUint(strings.Replace(r.Matcher, "cpu", "", 1), 10, 64)
		if err != nil {
			return false, nil, errors.Join(ErrInvalidRule, err)
		}
		idx = int64(idx_ + 1)
	}

	if idx < 0 || int(idx) >= len(ss) {
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("cpu index out of range: %d", idx))
	}
	s := ss[idx]
	switch t.ThresholdType {
	case ThresholdTypePercent:
		percent, err := s.UsedPercent()
		if err != nil {
			return false, nil, err
		}
		ok, err := t.True(percent)
		if err != nil {
			return false, nil, err
		}
		usedPercent, err := s.UsedPercent()
		if err != nil {
			return false, nil, err
		}
		key := "cpu"
		if idx > 0 {
			key = fmt.Sprintf("cpu%d", idx-1)
		}
		return ok, &PushPair{
			Key:   key,
			Value: fmt.Sprintf("%.2f%%", usedPercent),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid threshold type for cpu: %s", t.ThresholdType.Name()))
	}
}
func (r *Rule) shouldNotifyMemory(s *MemStatus, t *Threshold) (bool, *PushPair, error) {
	if s == nil {
		// utils.Warn("memory is not valid, skip this rule")
		return false, nil, nil
	}
	var size Size
	var percent float64
	switch r.Matcher {
	case "avail":
		size = s.Avail
		percent = float64(s.Avail) / float64(s.Total)
	case "free":
		size = s.Free
		percent = float64(s.Free) / float64(s.Total)
	case "used":
		size = s.Used
		percent = float64(s.Used) / float64(s.Total)
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid matcher: %s", r.Matcher))
	}

	switch t.ThresholdType {
	case ThresholdTypeSize:
		ok, err := t.True(size)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher + "of Memory",
			Value: size.String(),
		}, nil
	case ThresholdTypePercent:
		ok, err := t.True(percent)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher + "of Memory",
			Value: fmt.Sprintf("%.2f%%", percent*100),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid threshold type for memory: %s", t.ThresholdType.Name()))
	}
}
func (r *Rule) shouldNotifySwap(s *SwapStatus, t *Threshold) (bool, *PushPair, error) {
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
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher + "of Swap",
			Value: size.String(),
		}, nil
	case ThresholdTypePercent:
		ok, err := t.True(percent)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher + "of Swap",
			Value: fmt.Sprintf("%.2f%%", percent*100),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid threshold type for swap: %s", t.ThresholdType.Name()))
	}
}
func (r *Rule) shouldNotifyDisk(s []DiskStatus, t *Threshold) (bool, *PushPair, error) {
	if len(s) == 0 {
		// utils.Warn("disk is not valid, skip this rule")
		return false, nil, nil
	}
	var disk DiskStatus
	var have bool
	for _, d := range s {
		if d.MountPath == r.Matcher || d.Filesystem == r.Matcher {
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
		ok, err := t.True(disk.Used)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher,
			Value: disk.Used.String(),
		}, nil
	case ThresholdTypePercent:
		ok, err := t.True(disk.UsedPercent)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher,
			Value: fmt.Sprintf("%.2f%%", disk.UsedPercent),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid threshold type for disk: %s", t.ThresholdType.Name()))
	}
}
func (r *Rule) shouldNotifyNetwork(s []NetworkStatus, t *Threshold) (bool, *PushPair, error) {
	if len(s) == 0 {
		// utils.Warn("network is not valid, skip this rule")
		return false, nil, nil
	}

	var net NetworkStatus
	var have bool
	for _, n := range s {
		if strings.Contains(r.Matcher, n.Interface) {
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
		// 如果没有指定方向，则默认计算 出+入 流量
		in = true
		out = true
	}

	switch t.ThresholdType {
	case ThresholdTypeSpeed:
		speed := Size(0)
		if in {
			s, err := net.ReceiveSpeed()
			if err != nil {
				return false, nil, err
			}
			speed += s
		}
		if out {
			s, err := net.TransmitSpeed()
			if err != nil {
				return false, nil, err
			}
			speed += s
		}
		ok, err := t.True(speed)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher,
			Value: speed.String(),
		}, nil
	case ThresholdTypeSize:
		size := Size(0)
		if in {
			size += net.TimeSequence.New.Receive
		}
		if out {
			size += net.TimeSequence.New.Transmit
		}
		ok, err := t.True(size)
		if err != nil {
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher,
			Value: size.String(),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid threshold type for network: %s", t.ThresholdType.Name()))
	}
}

func (r *Rule) shouldNotifyTemperature(s []TemperatureStatus, t *Threshold) (bool, *PushPair, error) {
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
			return false, nil, err
		}
		return ok, &PushPair{
			Key:   r.Matcher,
			Value: fmt.Sprintf("%.2f°C", temp.Value),
		}, nil
	default:
		return false, nil, errors.Join(ErrInvalidRule, fmt.Errorf("invalid threshold type for temperature: %s", t.ThresholdType.Name()))
	}
}

type MonitorType string

const (
	MonitorTypeCPU         MonitorType = "cpu"
	MonitorTypeMemory                  = "mem"
	MonitorTypeSwap                    = "swap"
	MonitorTypeDisk                    = "disk"
	MonitorTypeNetwork                 = "net"
	MonitorTypeTemperature             = "temp"
)
