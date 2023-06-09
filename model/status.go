package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lollipopkit/gommon/term"
	"github.com/lollipopkit/gommon/util"
	"github.com/lollipopkit/server_box_monitor/res"
)

var (
	ErrNotReady           = errors.New("not ready")
	ErrRunShellFailed     = errors.New("run shell failed")
	ErrInvalidShellOutput = errors.New("invalid shell output")
)

var (
	Status = new(serverStatus)
)

type serverStatus struct {
	CPU         []oneCpuStatus
	Mem         *memStatus
	Swap        *swapStatus
	Disk        []diskStatus
	Network     []networkStatus
	Temperature []temperatureStatus
}

type temperatureStatus struct {
	Value float64
	Name  string
}

type cpuOneTimeStatus struct {
	Used  int
	Total int
}

type oneCpuStatus struct {
	Core int
	TimeSequence[cpuOneTimeStatus]
}

func (cs *oneCpuStatus) UsedPercent() (float64, error) {
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

type memStatus struct {
	Total Size
	Avail Size
	Free  Size
	Used  Size
}

type swapStatus struct {
	Total  Size
	Free   Size
	Used   Size
	Cached Size
}

type diskStatus struct {
	MountPath   string
	Filesystem  string
	Total       Size
	Used        Size
	Avail       Size
	UsedPercent float64
}

type networkOneTimeStatus struct {
	Transmit Size
	Receive  Size
}
type networkStatus struct {
	Interface string
	TimeSequence[networkOneTimeStatus]
}

func (ns *networkStatus) TransmitSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := float64(ns.TimeSequence.New.Transmit - ns.TimeSequence.Old.Transmit)
	return Size(diff / GetIntervalInSeconds()), nil
}
func (ns *networkStatus) ReceiveSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := float64(ns.TimeSequence.New.Receive - ns.TimeSequence.Old.Receive)
	return Size(diff / GetIntervalInSeconds()), nil
}

func RefreshStatus() error {
	output, _ := util.Execute("bash", res.ServerBoxShellPath)
	err := os.WriteFile(filepath.Join(res.ServerBoxDirPath, "shell_output.log"), []byte(output), 0644)
	if err != nil {
		term.Warn("[STATUS] write shell output log failed: %s", err)
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
		term.Warn("parse network status failed: %s", err)
	}
	err = parseCPUStatus(segments[2])
	if err != nil {
		term.Warn("parse cpu status failed: %s", err)
	}
	err = parseDiskStatus(segments[3])
	if err != nil {
		term.Warn("parse disk status failed: %s", err)
	}
	err = parseMemAndSwapStatus(segments[4])
	if err != nil {
		term.Warn("parse mem status failed: %s", err)
	}
	err = parseTemperatureStatus(segments[5], segments[6])
	if err != nil {
		term.Warn("parse temperature status failed: %s", err)
	}
	return nil
}

func initMem() {
	if Status.Mem == nil {
		Status.Mem = &memStatus{}
	}
}
func initSwap() {
	if Status.Swap == nil {
		Status.Swap = &swapStatus{}
	}
}

func parseMemAndSwapStatus(s string) error {
	initMem()
	initSwap()
	lines := strings.Split(s, "\n")
	for i := range lines {
		line := strings.TrimSpace(lines[i])

		value, err := strconv.ParseInt(strings.Fields(line)[1], 10, 64)
		if err != nil {
			return err
		}
		// KB -> B
		// because the unit of MemTotal/... is KB
		size := Size(value) * Size(programKilo)

		switch true {
		case strings.HasPrefix(line, "MemTotal:"):
			Status.Mem.Total = size
			fallthrough
		case strings.HasPrefix(line, "MemFree:"):
			Status.Mem.Free = size
			Status.Mem.Used = Status.Mem.Total - Status.Mem.Free
			fallthrough
		case strings.HasPrefix(line, "MemAvailable:"):
			Status.Mem.Avail = size
		case strings.HasPrefix(line, "SwapTotal:"):
			Status.Swap.Total = size
			fallthrough
		case strings.HasPrefix(line, "SwapFree:"):
			Status.Swap.Free = size
			Status.Swap.Used = Status.Swap.Total - Status.Swap.Free
			fallthrough
		case strings.HasPrefix(line, "SwapCached:"):
			Status.Swap.Cached = size
		}
	}
	return nil
}

func parseCPUStatus(s string) error {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	count := len(lines)
	if len(Status.CPU) != count {
		Status.CPU = make([]oneCpuStatus, count)
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
			Status.CPU[i].TimeSequence.Update(&cpuOneTimeStatus{
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
	if len(Status.Disk) != count {
		Status.Disk = make([]diskStatus, count)
	}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		fields := strings.Fields(line)
		if len(fields) != 6 {
			return errors.Join(ErrInvalidShellOutput, fmt.Errorf("invalid disk status: %s", line))
		}
		Status.Disk[i].MountPath = fields[5]
		Status.Disk[i].Filesystem = fields[0]
		total, err := ParseToSize(fields[1])
		if err != nil {
			return err
		}
		Status.Disk[i].Total = total
		used, err := ParseToSize(fields[2])
		if err != nil {
			return err
		}
		Status.Disk[i].Used = used
		avail, err := ParseToSize(fields[3])
		if err != nil {
			return err
		}
		Status.Disk[i].Avail = avail
		Status.Disk[i].UsedPercent = (float64(used) / float64(total)) * 100
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
	if len(Status.Temperature) != count {
		Status.Temperature = make([]temperatureStatus, count)
	}
	for i := range types {
		Status.Temperature[i].Name = strings.TrimSpace(types[i])
		value, err := strconv.ParseFloat(strings.TrimSpace(values[i]), 64)
		if err != nil {
			return err
		}
		Status.Temperature[i].Value = value / 1000
	}
	return nil
}

func parseNetworkStatus(s string) error {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	count := len(lines)
	if len(Status.Network) != count-2 {
		Status.Network = make([]networkStatus, count-2)
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
		Status.Network[idx].Interface = strings.TrimRight(fields[0], ":")
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
		Status.Network[idx].TimeSequence.Update(&networkOneTimeStatus{
			Receive:  receive,
			Transmit: transmit,
		})
	}
	return nil
}
