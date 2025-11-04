package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ethan-mdev/service-watch/internal/handlers"
	"github.com/ethan-mdev/service-watch/internal/monitor"
	"github.com/ethan-mdev/service-watch/internal/platform"
	"github.com/ethan-mdev/service-watch/internal/sse"
	"github.com/ethan-mdev/service-watch/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Initialize service manager
	svcMgr := platform.MakeServiceManager()

	// Initialize watchlist manager
	watchlistMgr := storage.NewJSONWatchlist("watchlist.json", svcMgr)

	// Initialize SSE broadcaster
	broadcaster := sse.NewBroadcaster()

	// Initialize service watcher
	go monitor.Start(context.Background(), watchlistMgr, svcMgr, broadcaster)

	// Create HTTP handlers
	svcHTTP := handlers.NewServiceHTTP(svcMgr)
	watchlistHTTP := handlers.NewWatchlistHTTP(watchlistMgr)
	eventsHTTP := handlers.NewEventsHTTP(broadcaster)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Mount routes
	r.Mount("/v1/services", svcHTTP.Routes())
	r.Mount("/v1/watchlist", watchlistHTTP.Routes())
	r.Get("/v1/events", eventsHTTP.Stream)

	// Start server
	addr := "127.0.0.1:8080"
	log.Printf("service-watch listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
