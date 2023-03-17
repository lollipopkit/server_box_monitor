package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	ErrNotReady           = errors.New("not ready")
	ErrRunShellFailed     = errors.New("run shell failed")
	ErrInvalidShellOutput = errors.New("invalid shell output")
)

var (
	status = new(Status)
)

func GetStatus() *Status {
	return status
}

type Status struct {
	CPU         []CPUStatus
	Mem         *MemStatus
	Swap        *SwapStatus
	Disk        []DiskStatus
	Network     []NetworkStatus
	Temperature []TemperatureStatus
}

type TemperatureStatus struct {
	Value float64
	Name  string
}

type TimeSequence[T CPUOneTimeStatus | NetworkOneTimeStatus] struct {
	Old *T
	New *T
}

func (ts *TimeSequence[T]) Update(t *T) {
	ts.Old = ts.New
	ts.New = t
}

type CPUOneTimeStatus struct {
	Used  int
	Total int
}
type CPUStatus struct {
	Core int
	TimeSequence[CPUOneTimeStatus]
}

func (cs *CPUStatus) UsedPercent() (float64, error) {
	if cs.TimeSequence.New == nil || cs.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	used := cs.TimeSequence.New.Used - cs.TimeSequence.Old.Used
	total := cs.TimeSequence.New.Total - cs.TimeSequence.Old.Total
	if total == 0 {
		return 0, ErrNotReady
	}
	return float64(used) / float64(total) * 100, nil
}

type MemStatus struct {
	Total Size
	Avail Size
	Free  Size
	Used  Size
}

type SwapStatus struct {
	Total  Size
	Free   Size
	Used   Size
	Cached Size
}

type DiskStatus struct {
	MountPath   string
	Filesystem  string
	Total       Size
	Used        Size
	Avail       Size
	UsedPercent float64
}

type NetworkOneTimeStatus struct {
	Transmit Size
	Receive  Size
}
type NetworkStatus struct {
	Interface string
	TimeSequence[NetworkOneTimeStatus]
}

func (ns *NetworkStatus) TransmitSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := ns.TimeSequence.New.Transmit - ns.TimeSequence.Old.Transmit
	return diff / Size(GetInterval()), nil
}
func (ns *NetworkStatus) ReceiveSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := ns.TimeSequence.New.Receive - ns.TimeSequence.Old.Receive
	return diff / Size(GetInterval()), nil
}

func RefreshStatus() error {
	output, _ := utils.Execute("bash", res.ServerBoxShellPath)
	err := os.WriteFile(filepath.Join(res.ServerBoxDirPath, "shell_output.log"), []byte(output), 0644)
	if err != nil {
		utils.Warn("[STATUS] write shell output log failed: %s", err)
	}
	return ParseStatus(output)
}

func ParseStatus(s string) error {
	segments := strings.Split(s, "SrvBox")
	for i := range segments {
		segments[i] = strings.TrimSpace(segments[i])
		segments[i] = strings.Trim(segments[i], "\n")
	}
	if len(segments) != 10 {
		return ErrInvalidShellOutput
	}
	err := parseNetworkStatus(segments[1])
	if err != nil {
		utils.Warn("parse network status failed: %s", err)
	}
	err = parseCPUStatus(segments[3])
	if err != nil {
		utils.Warn("parse cpu status failed: %s", err)
	}
	err = parseMemAndSwapStatus(segments[7])
	if err != nil {
		utils.Warn("parse mem status failed: %s", err)
	}
	err = parseTemperatureStatus(segments[8], segments[9])
	if err != nil {
		utils.Warn("parse temperature status failed: %s", err)
	}
	err = parseDiskStatus(segments[6])
	if err != nil {
		utils.Warn("parse disk status failed: %s", err)
	}
	return nil
}

func initMem() {
	if status.Mem == nil {
		status.Mem = &MemStatus{}
	}
}
func initSwap() {
	if status.Swap == nil {
		status.Swap = &SwapStatus{}
	}
}

func parseMemAndSwapStatus(s string) error {
	lines := strings.Split(s, "\n")
	for i := range lines {
		line := strings.TrimSpace(lines[i])

		value, err := strconv.ParseInt(strings.Fields(line)[1], 10, 64)
		if err != nil {
			return err
		}
		size := Size(value)

		switch true {
		case strings.HasPrefix(line, "MemTotal:"):
			initMem()
			status.Mem.Total = size
			fallthrough
		case strings.HasPrefix(line, "MemFree:"):
			initMem()
			status.Mem.Free = size
			status.Mem.Used = status.Mem.Total - status.Mem.Free
			fallthrough
		case strings.HasPrefix(line, "MemAvailable:"):
			initMem()
			status.Mem.Avail = size
		case strings.HasPrefix(line, "SwapTotal:"):
			initSwap()
			status.Swap.Total = size
			fallthrough
		case strings.HasPrefix(line, "SwapFree:"):
			initSwap()
			status.Swap.Free = size
			status.Swap.Used = status.Swap.Total - status.Swap.Free
			fallthrough
		case strings.HasPrefix(line, "SwapCached:"):
			initSwap()
			status.Swap.Cached = size
		}
	}
	return nil
}

func parseCPUStatus(s string) error {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	count := len(lines)
	if len(status.CPU) != count {
		status.CPU = make([]CPUStatus, count)
	}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			if len(fields) != 11 {
				return errors.Join(ErrInvalidShellOutput, fmt.Errorf("invalid cpu status: %s", line))
			}
			user, err := strconv.Atoi(fields[1])
			if err != nil {
				return err
			}
			sys, err := strconv.Atoi(fields[2])
			if err != nil {
				return err
			}
			total := 0
			for i := 2; i < 11; i++ {
				v, err := strconv.Atoi(fields[i])
				if err != nil {
					return err
				}
				total += v
			}
			status.CPU[i].TimeSequence.Update(&CPUOneTimeStatus{
				Used:  user + sys,
				Total: total,
			})
		}
	}
	return nil
}

func parseDiskStatus(s string) error {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	lines = lines[1:]
	count := len(lines)
	if len(status.Disk) != count {
		status.Disk = make([]DiskStatus, count)
	}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		fields := strings.Fields(line)
		if len(fields) != 6 {
			return errors.Join(ErrInvalidShellOutput, fmt.Errorf("invalid disk status: %s", line))
		}
		status.Disk[i].MountPath = fields[5]
		status.Disk[i].Filesystem = fields[0]
		total, err := ParseToSize(fields[1])
		if err != nil {
			return err
		}
		status.Disk[i].Total = total
		used, err := ParseToSize(fields[2])
		if err != nil {
			return err
		}
		status.Disk[i].Used = used
		avail, err := ParseToSize(fields[3])
		if err != nil {
			return err
		}
		status.Disk[i].Avail = avail
		status.Disk[i].UsedPercent = (float64(used) / float64(total)) * 100
	}
	return nil
}

func parseTemperatureStatus(s1, s2 string) error {
	if strings.Contains(s1, "/sys/class/thermal/thermal_zone*/type") {
		return nil
	}
	types := strings.Split(strings.TrimSpace(s1), "\n")
	values := strings.Split(strings.TrimSpace(s2), "\n")
	if len(types) != len(values) {
		return errors.Join(ErrInvalidShellOutput, fmt.Errorf("invalid temperature status: %s, %s", s1, s2))
	}
	count := len(types)
	if len(status.Temperature) != count {
		status.Temperature = make([]TemperatureStatus, count)
	}
	for i := range types {
		status.Temperature[i].Name = strings.TrimSpace(types[i])
		value, err := strconv.ParseFloat(strings.TrimSpace(values[i]), 64)
		if err != nil {
			return err
		}
		status.Temperature[i].Value = value / 1000
	}
	return nil
}

func parseNetworkStatus(s string) error {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	count := len(lines)
	if len(status.Network) != count-2 {
		status.Network = make([]NetworkStatus, count-2)
	}
	for i := range lines {
		if i < 2 {
			continue
		}
		line := strings.TrimSpace(lines[i])
		fields := strings.Fields(line)
		if len(fields) != 17 {
			return errors.Join(ErrInvalidShellOutput, fmt.Errorf("invalid network status: %s", line))
		}
		idx := i - 2
		status.Network[idx].Interface = strings.TrimRight(fields[0], ":")
		receiveBytes, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return err
		}
		receive := Size(receiveBytes)
		transmitBytes, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return err
		}
		transmit := Size(transmitBytes)
		status.Network[idx].TimeSequence.Update(&NetworkOneTimeStatus{
			Receive:  receive,
			Transmit: transmit,
		})
	}
	return nil
}
