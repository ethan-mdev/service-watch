# Service Watch

A lightweight Windows service monitoring agent with auto-restart capabilities, real-time metrics, and a REST API.

## Features

- **Service Management** - Start, stop, restart Windows services via REST API
- **Auto-Restart** - Monitor services and automatically restart on failure
- **Real-time Metrics** - CPU, memory, and uptime tracking for watched services
- **Live Updates** - Server-Sent Events (SSE) for real-time dashboard updates
- **Historical Data** - JSONL logging with automatic rotation and queryable API
- **Self-contained** - Single ~10MB binary, no external dependencies

## Quick Start

### Prerequisites

- Windows OS
- Administrator privileges (required for service management)

### Installation

1. Download the latest release or build from source:
```bash
go build -o service-watch.exe ./cmd/agent
```

2. Run as Administrator:
```bash
.\service-watch.exe
```

3. Access the API documentation:
```
http://localhost:8080/docs
```

## Usage

### Add a Service to Watchlist

```bash
curl -X POST http://localhost:8080/v1/watchlist \
  -H "Content-Type: application/json" \
  -d '{"serviceName":"Spooler","autoRestart":true}'
```

### Monitor Live Events

```bash
curl -N http://localhost:8080/v1/events
```

### Query Historical Metrics

```bash
# Last hour of metrics for a service
curl "http://localhost:8080/v1/metrics?event=metric_sample&service=Spooler&since=1h"
```

### Control Services

```bash
# Start service
curl http://localhost:8080/v1/services/Spooler/start

# Stop service
curl http://localhost:8080/v1/services/Spooler/stop

# Restart service
curl http://localhost:8080/v1/services/Spooler/restart
```

## API Endpoints

Full API documentation is available at `http://localhost:8080/docs` when running.

### Services
- `GET /v1/services` - List all services
- `GET /v1/services/{name}` - Get service details with metrics
- `GET /v1/services/{name}/start` - Start a service
- `GET /v1/services/{name}/stop` - Stop a service
- `GET /v1/services/{name}/restart` - Restart a service

### Watchlist
- `GET /v1/watchlist` - List monitored services
- `GET /v1/watchlist/{name}` - Get watchlist item
- `POST /v1/watchlist` - Add service to watchlist
- `PUT /v1/watchlist/{name}` - Update auto-restart setting
- `DELETE /v1/watchlist/{name}` - Remove from watchlist

### Metrics
- `GET /v1/metrics` - Query historical logs
  - Query params: `event`, `service`, `limit`, `since`

### Events (SSE)
- `GET /v1/events` - Real-time event stream

## Configuration

### Log Rotation

Logs are automatically rotated using the following defaults:
- **Max Size:** 10 MB
- **Max Backups:** 5 files
- **Max Age:** 7 days
- **Compression:** Enabled (gzip)

Logs are stored in `logs/events.jsonl` with rotated files compressed as `events-{timestamp}.jsonl.gz`.

### Watchlist Persistence

Watchlist configuration is stored in `watchlist.json` and persists across restarts.

## Architecture

```
service-watch/
├── cmd/agent/          # Main application entry point
├── internal/
│   ├── core/           # Domain types and interfaces
│   ├── handlers/       # HTTP handlers (services, watchlist, metrics, events)
│   ├── logger/         # Structured logging with SSE broadcasting
│   ├── monitor/        # Background service watcher
│   ├── platform/       # OS-specific service management (Windows/Linux)
│   ├── sse/            # Server-Sent Events broadcaster
│   ├── storage/        # Watchlist persistence
│   └── utils/          # HTTP utilities
├── static/             # API documentation
└── logs/               # Event logs (JSONL format)
```

## Event Types

The following events are broadcasted via SSE and logged to `logs/events.jsonl`:

- `app_started` - Application initialization
- `watcher_started` - Service monitor started
- `metric_sample` - Service metrics (CPU, memory, state)
- `restart_attempt` - Service restart initiated
- `restart_success` - Service restarted successfully
- `restart_failed` - Service restart failed
- `service_failed` - Service exceeded max restart attempts

## Platform Support

- ✅ **Windows** - Fully supported
- 🚧 **Linux** - Planned (systemd support)

## Development

### Build

```bash
go build -o service-watch.exe ./cmd/agent
```

### Run Tests

```bash
go test ./...
```

### Dependencies

- [chi](https://github.com/go-chi/chi) - HTTP router
- [gopsutil](https://github.com/shirou/gopsutil) - Process metrics
- [lumberjack](https://github.com/natefinch/lumberjack) - Log rotation
- [golang.org/x/sys](https://golang.org/x/sys) - Windows service control

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Svelte dashboard with real-time graphs
- [ ] Discord webhook integration for alerts
- [ ] Linux/systemd support
- [ ] Authentication for remote access
