package model

import (
	"errors"
	"strings"

	"github.com/lollipopkit/server_box_monitor/res"
	"github.com/lollipopkit/server_box_monitor/utils"
)

var (
	ErrNotReady = errors.New("not ready")
	ErrRunShellFailed = errors.New("run shell failed")
	ErrInvalidShellOutput = errors.New("invalid shell output")
)

var (
	status *Status
)
func GetStatus() *Status {
	return status
}

type Status struct {
	CPU     []CPUStatus
	Mem     *MemStatus
	Swap    *SwapStatus
	Disk    []DiskStatus
	Network []NetworkStatus
	Temperature []TemperatureStatus
}

type TemperatureStatus struct {
	Value float64
	Name string
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
	Used int
	Total int
}
type CPUStatus struct {
	Core        int
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

type NetworkOneTimeStatus struct {
	Transmit Size
	Receive Size
}
type NetworkStatus struct {
	Interface     string
	TimeSequence[NetworkOneTimeStatus]
	TransmitAmount Size
	ReceiveAmount Size
}
func (ns *NetworkStatus) TransmitSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := ns.TimeSequence.New.Transmit - ns.TimeSequence.Old.Transmit
	return diff / Size(Config.GetRunInterval()), nil
}
func (ns *NetworkStatus) ReceiveSpeed() (Size, error) {
	if ns.TimeSequence.New == nil || ns.TimeSequence.Old == nil {
		return 0, ErrNotReady
	}
	diff := ns.TimeSequence.New.Receive - ns.TimeSequence.Old.Receive
	return diff / Size(Config.GetRunInterval()), nil
}

func RefreshStatus() error {
	stdout, stderr, err := utils.Execute("bash", res.ServerBoxShellPath)
	if err != nil {
		utils.Warn("run shell failed: %s, %s", err, stderr)
	}
	return ParseStatus(stdout)
}

func ParseStatus(s string) error {
	segments := strings.Split(s, "SrvBox")
	for i, segment := range segments {
		segments[i] = strings.TrimSpace(segment)
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
	return ErrRunShellFailed
}

func parseMemAndSwapStatus(s string) error {
	lines := strings.Split(s, "\n")
	return nil
}
