package handlers

import (
	"fmt"
	"log"
	"net/http"

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
		Channel: make(chan sse.Event, 10),
	}

	// Register client
	h.Broadcaster.RegisterClient(client)
	defer h.Broadcaster.UnregisterClient(client)

	log.Printf("SSE: Client connected: %v", r.RemoteAddr)

	// Send initial connection event
	fmt.Fprint(w, sse.FormatEvent(sse.Event{
		Type: "connected",
		Data: map[string]interface{}{
			"message": "Connected to service-watch events",
		},
	}))
	w.(http.Flusher).Flush()

	// Stream events
	for {
		select {
		case <-r.Context().Done():
			return
		case event := <-client.Channel:
			fmt.Fprint(w, sse.FormatEvent(event))
			w.(http.Flusher).Flush()
		}
	}
}
