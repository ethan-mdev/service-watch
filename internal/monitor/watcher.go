package monitor

import (
	"context"
	"log"
	"time"

	"github.com/ethan-mdev/service-watch/internal/core"
)

// Start begins monitoring watchlist items and auto-restarting services.
func Start(ctx context.Context, watchlistMgr core.WatchlistManager, svcMgr core.ServiceManager) {
	log.Println("Starting service monitor...")

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Check immediately on startup
	checkServices(ctx, watchlistMgr, svcMgr)

	for {
		select {
		case <-ctx.Done():
			log.Println("Service watcher stopped")
			return
		case <-ticker.C:
			checkServices(ctx, watchlistMgr, svcMgr)
		}
	}
}

func checkServices(ctx context.Context, watchlistMgr core.WatchlistManager, svcMgr core.ServiceManager) {
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
				watchlistMgr.Update(ctx, item.ServiceName, false)
				continue
			}

			log.Printf("Watcher: service %s is %s, attempting restart...",
				item.ServiceName, item.Service.State)

			if err := svcMgr.Start(ctx, item.ServiceName); err != nil {
				item.FailCount++
				log.Printf("Watcher: failed to restart %s: %v", item.ServiceName, err)
			} else {
				log.Printf("Watcher: successfully restarted %s", item.ServiceName)
				watchlistMgr.IncrementRestartCount(ctx, item.ServiceName)
			}
		}
	}
}
