package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethan-mdev/service-watch/internal/utils"
	"github.com/go-chi/chi/v5"
)

type MetricsHTTP struct {
	LogPath string
}

func NewMetricsHTTP(logPath string) *MetricsHTTP {
	return &MetricsHTTP{LogPath: logPath}
}

func (h *MetricsHTTP) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.query)
	return r
}

// query handles filtering and returning log entries
func (h *MetricsHTTP) query(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	eventType := r.URL.Query().Get("event")     // e.g., "watcher_started, host_resources, service_status"
	serviceName := r.URL.Query().Get("service") // e.g., "Steam Client Service"
	limitStr := r.URL.Query().Get("limit")      // e.g., "100"
	sinceStr := r.URL.Query().Get("since")      // e.g., "1h" or RFC3339 timestamp

	limit := 100 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// Parse "since" parameter
	var sinceTime time.Time
	if sinceStr != "" {
		// Try parsing as duration (e.g., "1h", "30m")
		if d, err := time.ParseDuration(sinceStr); err == nil {
			sinceTime = time.Now().Add(-d)
		} else if t, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			// Try parsing as timestamp
			sinceTime = t
		}
	}

	// Read and filter logs
	file, err := os.Open(h.LogPath)
	if err != nil {
		utils.RespondWithError(w, 500, "failed to open log file", err)
		return
	}
	defer file.Close()

	var results []map[string]interface{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var entry map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		// Filter by event type
		if eventType != "" && entry["event"] != eventType {
			continue
		}

		// Filter by service name
		if serviceName != "" {
			if data, ok := entry["data"].(map[string]interface{}); ok {
				if data["serviceName"] != serviceName {
					continue
				}
			} else {
				continue
			}
		}

		// Filter by time
		if !sinceTime.IsZero() {
			if timeStr, ok := entry["time"].(string); ok {
				if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
					if t.Before(sinceTime) {
						continue
					}
				}
			}
		}

		results = append(results, entry)

		// Limit results
		if len(results) >= limit {
			break
		}
	}

	utils.RespondWithJSON(w, 200, map[string]interface{}{
		"count": len(results),
		"items": results,
	})
}
