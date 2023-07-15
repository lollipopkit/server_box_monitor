package web

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/lollipopkit/server_box_monitor/model"
)

func Status(c echo.Context) error {
	s := model.Status
	var cpu any = 0.0
	if len(s.CPU) > 0 {
		var err error
		cpu, err = s.CPU[0].UsedPercent()
		if err != nil {
			cpu = err
		}
	}
	mem := ""
	if s.Mem != nil {
		mem = fmt.Sprintf("%s / %s", s.Mem.Used.String(), s.Mem.Total.String())
	}
	net := ""
	if len(s.Network) > 0 {
		all := model.AllNetworkStatus(s.Network)
		rx, _ := all.ReceiveSpeed()
		tx, _ := all.TransmitSpeed()
		net = fmt.Sprintf("%s / %s", rx.String(), tx.String())
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

	status := map[string]any{
		"cpu":  cpu,
		"mem":  mem,
		"net":  net,
		"disk": disk,
	}
	return c.JSON(200, status)
}
