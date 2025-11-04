package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ethan-mdev/service-watch/internal/core"
)

type jsonWatchlist struct {
	mutex      sync.RWMutex
	filepath   string
	svcManager core.ServiceManager
	items      map[string]*core.WatchlistItem
}

func NewJSONWatchlist(filepath string, svcManager core.ServiceManager) core.WatchlistManager {
	watchList := &jsonWatchlist{
		filepath:   filepath,
		svcManager: svcManager,
		items:      make(map[string]*core.WatchlistItem),
	}
	watchList.load()
	return watchList
}

// List implements core.WatchlistManager.
func (j *jsonWatchlist) List(ctx context.Context) ([]core.WatchlistItem, error) {
	j.mutex.RLock()
	defer j.mutex.RUnlock()

	items := make([]core.WatchlistItem, 0, len(j.items))
	for _, item := range j.items {
		result := *item // copy
		// Populate current service state
		if svc, err := j.svcManager.Get(ctx, item.ServiceName); err == nil {
			result.Service = &svc
		}
		items = append(items, result)
	}
	return items, nil
}

// Get implements core.WatchlistManager.
func (j *jsonWatchlist) Get(ctx context.Context, serviceName string) (core.WatchlistItem, error) {
	j.mutex.RLock()
	item, exists := j.items[serviceName]
	j.mutex.RUnlock()

	if !exists {
		return core.WatchlistItem{}, fmt.Errorf("service not in watchlist: %s", serviceName)
	}

	result := *item
	// Populate current service state
	if svc, err := j.svcManager.Get(ctx, serviceName); err == nil {
		result.Service = &svc
	}
	return result, nil
}

// Add implements core.WatchlistManager.
func (j *jsonWatchlist) Add(ctx context.Context, serviceName string, autoRestart bool) error {
	// Verify service exists
	if _, err := j.svcManager.Get(ctx, serviceName); err != nil {
		return fmt.Errorf("service not found: %s", serviceName)
	}

	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, exists := j.items[serviceName]; exists {
		return fmt.Errorf("service already in watchlist: %s", serviceName)
	}

	j.items[serviceName] = &core.WatchlistItem{
		ServiceName:  serviceName,
		AutoRestart:  autoRestart,
		RestartCount: 0,
		FailCount:    0,
	}
	return j.save()
}

// Remove implements core.WatchlistManager.
func (j *jsonWatchlist) Remove(ctx context.Context, serviceName string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	if _, exists := j.items[serviceName]; !exists {
		return fmt.Errorf("service not in watchlist: %s", serviceName)
	}

	delete(j.items, serviceName)
	return j.save()
}

// Update implements core.WatchlistManager.
func (j *jsonWatchlist) Update(ctx context.Context, serviceName string, autoRestart bool) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	item, exists := j.items[serviceName]
	if !exists {
		return fmt.Errorf("service not in watchlist: %s", serviceName)
	}
	if autoRestart {
		item.FailCount = 0
	}
	item.AutoRestart = autoRestart
	return j.save()
}

func (j *jsonWatchlist) load() error {
	data, err := os.ReadFile(j.filepath)
	if os.IsNotExist(err) {
		return nil // fresh start
	}
	if err != nil {
		return err
	}

	var items []core.WatchlistItem
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	for i := range items {
		j.items[items[i].ServiceName] = &items[i]
	}
	return nil
}

func (j *jsonWatchlist) save() error {
	items := make([]core.WatchlistItem, 0, len(j.items))
	for _, item := range j.items {
		// Don't save the embedded service data, just the watchlist config
		items = append(items, core.WatchlistItem{
			ServiceName:  item.ServiceName,
			AutoRestart:  item.AutoRestart,
			RestartCount: item.RestartCount,
			LastRestart:  item.LastRestart,
			FailCount:    item.FailCount,
		})
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(j.filepath, data, 0644)
}

// UpdateRestartInfo implements core.WatchlistManager.
func (j *jsonWatchlist) IncrementRestartCount(ctx context.Context, serviceName string) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	item, exists := j.items[serviceName]
	if !exists {
		return fmt.Errorf("service not in watchlist: %s", serviceName)
	}

	item.RestartCount++
	item.LastRestart = time.Now().Format(time.RFC3339)
	return j.save()
}
