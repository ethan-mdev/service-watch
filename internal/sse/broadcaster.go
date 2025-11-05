package sse

import (
	"sync"

	"github.com/ethan-mdev/service-watch/internal/core"
)

// Represents an SSE client.
type Client struct {
	Channel chan core.Event
}

// Manages SSE clients and broadcasts events to them.
type Broadcaster struct {
	clients map[*Client]bool
	mutex   sync.RWMutex
}

// Creates a new broadcaster.
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[*Client]bool),
	}
}

// Registers a new client to receive events.
func (b *Broadcaster) RegisterClient(client *Client) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.clients[client] = true
}

// Unregisters a client and closes its channel.
func (b *Broadcaster) UnregisterClient(client *Client) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if _, ok := b.clients[client]; ok {
		delete(b.clients, client)
		close(client.Channel)
	}
}

// Broadcasts an event to all registered clients.
func (b *Broadcaster) Broadcast(event core.Event) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	for client := range b.clients {
		select {
		case client.Channel <- event:
		default:
		}
	}
}
