package model

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/lollipopkit/gommon/log"
	"github.com/lollipopkit/gommon/sys"
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

type networkIface interface {
	TransmitSpeed() (Size, error)
	ReceiveSpeed() (Size, error)
	Transmit() Size
	Receive() Size
}

type networkStatus struct {
	Interface string
	TimeSequence[networkOneTimeStatus]
}

func (ns networkStatus) TransmitSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := float64(ns.TimeSequence.New.Transmit - ns.TimeSequence.Old.Transmit)
	return Size(diff / CheckInterval.Seconds()), nil
}
func (ns networkStatus) ReceiveSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := float64(ns.TimeSequence.New.Receive - ns.TimeSequence.Old.Receive)
	return Size(diff / CheckInterval.Seconds()), nil
}
func (ns networkStatus) Transmit() Size {
	return ns.TimeSequence.New.Transmit
}
func (ns networkStatus) Receive() Size {
	return ns.TimeSequence.New.Receive
}

type AllNetworkStatus []networkStatus

func (nss AllNetworkStatus) TransmitSpeed() (Size, error) {
	var sum float64
	for _, ns := range nss {
		speed, err := ns.TransmitSpeed()
		if err != nil {
			return 0, err
		}
		sum += float64(speed)
	}
	return Size(sum), nil
}
func (nss AllNetworkStatus) ReceiveSpeed() (Size, error) {
	var sum float64
	for _, ns := range nss {
		speed, err := ns.ReceiveSpeed()
		if err != nil {
			return 0, err
		}
		sum += float64(speed)
	}
	return Size(sum), nil
}
func (nss AllNetworkStatus) Transmit() Size {
	var sum float64
	for _, ns := range nss {
		sum += float64(ns.Transmit())
	}
	return Size(sum)
}
func (nss AllNetworkStatus) Receive() Size {
	var sum float64
	for _, ns := range nss {
		sum += float64(ns.Receive())
	}
	return Size(sum)
}

func RefreshStatus() error {
	output, _ := sys.Execute("sh", res.ServerBoxShellPath)
	err := os.WriteFile(filepath.Join(res.ServerBoxDirPath, "shell_output.log"), []byte(output), 0644)
	if err != nil {
		log.Warn("[STATUS] write shell output log failed: %s", err)
	}
	return ParseStatus(output)
}

func ParseStatus(s string) error {
	segments := strings.Split(s, "SrvBox")
	for i := range segments {
		segments[i] = strings.TrimSpace(segments[i])
		segments[i] = strings.Trim(segments[i], "\n")
	}
	if len(segments) != 7 {
		return errors.Join(ErrInvalidShellOutput, fmt.Errorf("expect 7 segments, but got %d", len(segments)))
	}
	err := ParseNetworkStatus(segments[1])
	if err != nil {
		log.Warn("parse network status failed: %s", err)
	}
	err = ParseCPUStatus(segments[2])
	if err != nil {
		log.Warn("parse cpu status failed: %s", err)
	}
	err = ParseDiskStatus(segments[3])
	if err != nil {
		log.Warn("parse disk status failed: %s", err)
	}
	err = ParseMemAndSwapStatus(segments[4])
	if err != nil {
		log.Warn("parse mem status failed: %s", err)
	}
	err = ParseTemperatureStatus(segments[5], segments[6])
	if err != nil {
		log.Warn("parse temperature status failed: %s", err)
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

func ParseMemAndSwapStatus(s string) error {
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
		case strings.HasPrefix(line, "MemFree:"):
			Status.Mem.Free = size
		case strings.HasPrefix(line, "MemAvailable:"):
			Status.Mem.Avail = size
		case strings.HasPrefix(line, "SwapTotal:"):
			Status.Swap.Total = size
		case strings.HasPrefix(line, "SwapFree:"):
			Status.Swap.Free = size
		case strings.HasPrefix(line, "SwapCached:"):
			Status.Swap.Used = size
		}
	}
	Status.Mem.Used = Status.Mem.Total - Status.Mem.Avail
	return nil
}

func ParseCPUStatus(s string) error {
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
			idle, err := strconv.Atoi(fields[4])
			if err != nil {
				return err
			}
			total := 0
			for i := 1; i < 8; i++ {
				v, err := strconv.Atoi(fields[i])
				if err != nil {
					return err
				}
				total += v
			}
			Status.CPU[i].TimeSequence.Update(&cpuOneTimeStatus{
				Used:  total - idle,
				Total: total,
			})
		}
	}
	return nil
}

func ParseDiskStatus(s string) error {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	lines = lines[1:]
	count := len(lines)
	if len(Status.Disk) != count {
		Status.Disk = make([]diskStatus, count)
	}
	failIdx := sort.IntSlice{}
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		fields := strings.Fields(line)
		if len(fields) != 6 {
			failIdx = append(failIdx, i)
			continue
		}
		Status.Disk[i].MountPath = fields[5]
		Status.Disk[i].Filesystem = fields[0]
		total, err := ParseToSize(fields[1])
		if err != nil {
			failIdx = append(failIdx, i)
			continue
		}
		Status.Disk[i].Total = total
		used, err := ParseToSize(fields[2])
		if err != nil {
			failIdx = append(failIdx, i)
			continue
		}
		Status.Disk[i].Used = used
		avail, err := ParseToSize(fields[3])
		if err != nil {
			failIdx = append(failIdx, i)
			continue
		}
		Status.Disk[i].Avail = avail
		Status.Disk[i].UsedPercent = (float64(used) / float64(total)) * 100
	}
	sort.Sort(sort.Reverse(failIdx))
	for _, idx := range failIdx {
		Status.Disk = append(Status.Disk[:idx], Status.Disk[idx+1:]...)
	}
	return nil
}

func ParseTemperatureStatus(s1, s2 string) error {
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

func ParseNetworkStatus(s string) error {
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
