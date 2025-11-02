package system

import (
	"time"

	"github.com/ethan-mdev/service-watch/internal/services"
	"github.com/shirou/gopsutil/v4/host"
)

type System struct {
	Start    time.Time
	HostInfo *host.InfoStat
	Services services.ServiceLister
}

func InitSystem() (*System, error) {
	info, err := host.Info()
	if err != nil {
		return nil, err
	}
	return &System{
		Start:    time.Now(),
		HostInfo: info,
		Services: services.NewServiceLister(), // selected by build tags
	}, nil
}
