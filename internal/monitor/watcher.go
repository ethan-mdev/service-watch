package monitor

import (
	"context"
	"time"

	"github.com/ethan-mdev/service-watch/internal/core"
	"github.com/ethan-mdev/service-watch/internal/logger"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// Start begins monitoring watchlist items and auto-restarting services.
func Start(ctx context.Context, watchlistMgr core.WatchlistManager, svcMgr core.ServiceManager, log *logger.Logger) {
	log.Info("watcher_started", map[string]interface{}{
		"interval": "10s",
	})

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Check immediately on startup
	checkServices(ctx, watchlistMgr, svcMgr, log)

	for {
		select {
		case <-ctx.Done():
			log.Info("watcher_stopped", nil)
			return
		case <-ticker.C:
			checkServices(ctx, watchlistMgr, svcMgr, log)
		}
	}
}

func checkServices(ctx context.Context, watchlistMgr core.WatchlistManager, svcMgr core.ServiceManager, log *logger.Logger) {
	items, err := watchlistMgr.List(ctx)
	if err != nil {
		log.Error("watcher_list_failed", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	checkHostResources(log)

	for _, item := range items {

		if item.Service != nil {
			log.Info("service_status", map[string]interface{}{
				"serviceName": item.ServiceName,
				"state":       item.Service.State,
				"cpuPercent":  item.Service.CPUPercent,
				"memoryMB":    item.Service.MemoryMB,
				"uptimeSec":   item.Service.UptimeSeconds,
				"pid":         item.Service.PID,
			})
		}

		if !item.AutoRestart || item.Service == nil {
			continue
		}

		if item.Service.State != "running" {
			if item.FailCount >= 3 {
				log.Error("service_failed", map[string]interface{}{
					"serviceName": item.ServiceName,
					"failCount":   item.FailCount,
					"message":     "Exceeded max restart attempts",
				})
				watchlistMgr.Update(ctx, item.ServiceName, false)
				continue
			}

			log.Info("restart_attempt", map[string]interface{}{
				"serviceName": item.ServiceName,
				"state":       item.Service.State,
			})

			if err := svcMgr.Start(ctx, item.ServiceName); err != nil {
				log.Error("restart_failed", map[string]interface{}{
					"serviceName": item.ServiceName,
					"error":       err.Error(),
					"failCount":   item.FailCount + 1,
				})
			} else {
				watchlistMgr.IncrementRestartCount(ctx, item.ServiceName)
				log.Info("restart_success", map[string]interface{}{
					"serviceName":  item.ServiceName,
					"restartCount": item.RestartCount + 1,
				})
			}
		}
	}
}

func checkHostResources(log *logger.Logger) {
	mem, _ := mem.VirtualMemory()
	cpuPercents, _ := cpu.Percent(time.Second, false)
	log.Info("host_resources", map[string]interface{}{
		"cpuPercent":  cpuPercents[0],
		"totalMB":     mem.Total / 1024 / 1024,
		"usedMB":      mem.Used / 1024 / 1024,
		"usedPercent": mem.UsedPercent,
	})
}
