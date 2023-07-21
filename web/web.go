package web

import (
	"fmt"

	"github.com/labstack/echo/v4"
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
	disk := ""
	if len(s.Disk) > 0 {
		d := s.Disk[0]
		for _, v := range s.Disk {
			if v.MountPath == "/" {
				d = v
				break
			}
		}
		disk = fmt.Sprintf("%s / %s", d.Used.String(), d.Total.String())
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
