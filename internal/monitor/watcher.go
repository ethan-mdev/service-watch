package monitor

import (
	"context"
	"log"
	"time"

	"github.com/ethan-mdev/service-watch/internal/core"
	"github.com/ethan-mdev/service-watch/internal/sse"
)

// Start begins monitoring watchlist items and auto-restarting services.
func Start(ctx context.Context, watchlistMgr core.WatchlistManager, svcMgr core.ServiceManager, broadcaster *sse.Broadcaster) {
	log.Println("Starting service monitor...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Check immediately on startup
	checkServices(ctx, watchlistMgr, svcMgr, broadcaster)

	for {
		select {
		case <-ctx.Done():
			log.Println("Service watcher stopped")
			return
		case <-ticker.C:
			checkServices(ctx, watchlistMgr, svcMgr, broadcaster)
		}
	}
}

func checkServices(ctx context.Context, watchlistMgr core.WatchlistManager, svcMgr core.ServiceManager, broadcaster *sse.Broadcaster) {
	items, err := watchlistMgr.List(ctx)
	if err != nil {
		log.Printf("Watcher: failed to list watchlist: %v", err)
		return
	}

	for _, item := range items {
		// Only restart if:
		// 1. AutoRestart is enabled
		// 2. Service info was retrieved successfully
		// 3. Service is not running
		if !item.AutoRestart || item.Service == nil {
			continue
		}

		if item.Service.State != "running" {
			if item.FailCount >= 3 {
				log.Printf("Watcher: service %s has failed %d times, skipping restart", item.ServiceName, item.FailCount)

				broadcaster.Broadcast(sse.Event{
					Type: "service_failed",
					Data: map[string]interface{}{
						"service_name": item.ServiceName,
						"fail_count":   item.FailCount,
						"message":      "Service has failed multiple times, manual intervention required",
						"timestamp":    time.Now(),
					},
				})
				watchlistMgr.Update(ctx, item.ServiceName, false)
				continue
			}

			log.Printf("Watcher: service %s is %s, attempting restart...", item.ServiceName, item.Service.State)

			broadcaster.Broadcast(sse.Event{
				Type: "service_restarting",
				Data: map[string]interface{}{
					"service_name": item.ServiceName,
					"timestamp":    time.Now(),
				},
			})

			if err := svcMgr.Start(ctx, item.ServiceName); err != nil {
				item.FailCount++
				log.Printf("Watcher: failed to restart %s: %v", item.ServiceName, err)

				broadcaster.Broadcast(sse.Event{
					Type: "service_restart_failed",
					Data: map[string]interface{}{
						"service_name": item.ServiceName,
						"fail_count":   item.FailCount,
						"message":      "Service failed to restart",
						"timestamp":    time.Now(),
					},
				})
			} else {
				log.Printf("Watcher: successfully restarted %s", item.ServiceName)
				watchlistMgr.IncrementRestartCount(ctx, item.ServiceName)

				broadcaster.Broadcast(sse.Event{
					Type: "service_restart_success",
					Data: map[string]interface{}{
						"service_name":  item.ServiceName,
						"restart_count": item.RestartCount,
						"timestamp":     time.Now(),
					},
				})
			}
		}
	}
}
