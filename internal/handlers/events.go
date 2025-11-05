package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethan-mdev/service-watch/internal/core"
	"github.com/ethan-mdev/service-watch/internal/sse"
)

type EventsHTTP struct {
	Broadcaster *sse.Broadcaster
}

func NewEventsHTTP(broadcaster *sse.Broadcaster) *EventsHTTP {
	return &EventsHTTP{Broadcaster: broadcaster}
}

// Stream handles SSE connections.
func (h *EventsHTTP) Stream(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create client
	client := &sse.Client{
		Channel: make(chan core.Event, 10),
	}

	// Register client
	h.Broadcaster.RegisterClient(client)
	defer h.Broadcaster.UnregisterClient(client)

	log.Printf("SSE: Client connected: %v", r.RemoteAddr)

	// Stream events
	for {
		select {
		case <-r.Context().Done():
			log.Printf("SSE: Client disconnected: %v", r.RemoteAddr)
			return
		case event := <-client.Channel:
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, toJSON(event.Data))
			w.(http.Flusher).Flush()
		}
	}
}

func toJSON(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}
