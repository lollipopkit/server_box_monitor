package web

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/gommon/util"
	"github.com/lollipopkit/server_box_monitor/model"
)

func Status(c echo.Context) error {
	s := model.Status
	cpu := ""
	if len(s.CPU) > 0 {
		cpu_, _ := s.CPU[0].UsedPercent()
		cpu = fmt.Sprintf("%.1f%%", cpu_)
	}
	mem := ""
	if s.Mem != nil {
		mem = fmt.Sprintf("%s / %s", s.Mem.Used.String(), s.Mem.Total.String())
	}
	net := ""
	if len(s.Network) > 0 {
		all := model.AllNetworkStatus(s.Network)
		rx := all.Receive().String()
		tx := all.Transmit().String()
		net = fmt.Sprintf("%s / %s", rx, tx)
	}
	diskUsed := model.Size(0)
	diskTotal := model.Size(0)
	diskDevs := []string{}
	for _, v := range s.Disk {
		if !strings.HasPrefix(v.Filesystem, "/dev") {
			continue
		}
		if util.Contains(diskDevs, v.Filesystem) {
			continue
		}
		diskDevs = append(diskDevs, v.Filesystem)
		diskUsed += v.Used
		diskTotal += v.Total
	}
	disk := ""
	if diskTotal > 0 {
		disk = fmt.Sprintf("%s / %s", diskUsed.String(), diskTotal.String())
	}
	status := map[string]string{
		"name": model.Config.Name,
		"cpu":  cpu,
		"mem":  mem,
		"net":  net,
		"disk": disk,
	}
	return ok(c, status)
}
