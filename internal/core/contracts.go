package core

import "context"

// ServiceManager abstracts OS service control.
type ServiceManager interface {
	// Lists all services.
	List(ctx context.Context) ([]Service, error)
	// Gets a specific service by name with full details including resource usage.
	Get(ctx context.Context, name string) (Service, error)
	// Starts a service by name.
	Start(ctx context.Context, name string) error
	// Stops a service by name.
	Stop(ctx context.Context, name string) error
	// Restarts a service by name.
	Restart(ctx context.Context, name string) error
}

// WatchlistManager abstracts watchlist management.
type WatchlistManager interface {
	// Lists all watchlist items with current service details populated.
	List(ctx context.Context) ([]WatchlistItem, error)
	// Gets a specific watchlist item by service name with current service details.
	Get(ctx context.Context, name string) (WatchlistItem, error)
	// Adds a service to the watchlist with auto-restart configuration.
	Add(ctx context.Context, name string, autoRestart bool) error
	// Removes a service from the watchlist.
	Remove(ctx context.Context, name string) error
	// Updates the auto-restart setting for a watchlist item.
	Update(ctx context.Context, name string, autoRestart bool) error
}
