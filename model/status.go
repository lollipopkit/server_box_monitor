package model

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

func GetStatus() (*Status, error) {
	// TODO
	return nil, ErrCompareFailed
}
