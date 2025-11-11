package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os/exec"

	"github.com/ethan-mdev/service-watch/internal/handlers"
	"github.com/ethan-mdev/service-watch/internal/logger"
	"github.com/ethan-mdev/service-watch/internal/monitor"
	"github.com/ethan-mdev/service-watch/internal/platform"
	"github.com/ethan-mdev/service-watch/internal/sse"
	"github.com/ethan-mdev/service-watch/internal/storage"
	"github.com/getlantern/systray"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed dist
var webFS embed.FS

//go:embed icon.ico
var iconData []byte

func main() {
	// Start the server in a goroutine so it doesn't block
	go startServer()

	// Now run the system tray (this blocks)
	systray.Run(onTrayReady, onTrayExit)
}

func startServer() {
	// Initialize SSE broadcaster
	broadcaster := sse.NewBroadcaster()

	// Initialize logger with broadcaster
	appLogger, err := logger.Start("logs/events.jsonl", broadcaster)
	if err != nil {
		// Can't use log package in GUI app, so just panic on critical errors
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer appLogger.Close()

	// Initialize service manager
	svcMgr := platform.MakeServiceManager()

	// Initialize watchlist manager
	watchlistMgr := storage.NewJSONWatchlist("watchlist.json", svcMgr)

	// Initialize service watcher with logger
	go monitor.Start(context.Background(), watchlistMgr, svcMgr, appLogger)

	// Create HTTP handlers
	svcHTTP := handlers.NewServiceHTTP(svcMgr)
	watchlistHTTP := handlers.NewWatchlistHTTP(watchlistMgr)
	eventsHTTP := handlers.NewEventsHTTP(broadcaster)
	metricsHTTP := handlers.NewMetricsHTTP("logs/events.jsonl")

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Mount API routes
	r.Mount("/v1/services", svcHTTP.Routes())
	r.Mount("/v1/watchlist", watchlistHTTP.Routes())
	r.Mount("/v1/metrics", metricsHTTP.Routes())
	r.Get("/v1/events", eventsHTTP.Stream)

	// Serve API docs at /docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/docs.html")
	})

	// Serve embedded web app
	distFS, err := fs.Sub(webFS, "dist/")
	if err != nil {
		// Log error
		appLogger.Error("embed_error", map[string]interface{}{
			"error":    err.Error(),
			"fallback": "serving from disk",
		})
		r.Handle("/*", http.FileServer(http.Dir("dist/")))
	} else {
		r.Handle("/*", http.FileServer(http.FS(distFS)))
	}

	// Log server startup
	addr := "127.0.0.1:8080"
	appLogger.Info("server_starting", map[string]interface{}{
		"address": addr,
	})

	// Start server - use panic instead of log.Fatal for GUI apps
	if err := http.ListenAndServe(addr, r); err != nil {
		appLogger.Error("server_error", map[string]interface{}{
			"error": err.Error(),
		})
		panic(fmt.Sprintf("Server failed: %v", err))
	}
}

func onTrayReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Service Watch")
	systray.SetTooltip("Service Watch - Running on port 8080")

	mOpen := systray.AddMenuItem("Open Dashboard", "Open in browser")
	mLogs := systray.AddMenuItem("Open Logs Folder", "Open logs directory")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Exit", "Exit Service Watch")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openBrowser("http://127.0.0.1:8080")
			case <-mLogs.ClickedCh:
				openLogsFolder()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onTrayExit() {
	// Cleanup
}

func openBrowser(url string) {
	exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}

func openLogsFolder() {
	// Open the logs folder in current directory
	exec.Command("explorer", "logs").Start()
}
