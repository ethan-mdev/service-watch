# Service Watch

A standalone Windows service monitoring application with auto-restart capabilities, real-time web dashboard, and system tray integration.

## Features

- **Web Dashboard** - Modern real-time interface with live charts and service management
- **System Tray Integration** - Runs quietly in background, accessible via system tray icon
- **Service Management** - Start, stop, restart Windows services with one click
- **Auto-Restart** - Monitor services and automatically restart on failure
- **Real-time Metrics** - CPU, memory, and uptime tracking with live updating charts
- **Live Log Streaming** - View service events in real-time through the web interface
- **Historical Data** - Queryable event logs with filtering and search capabilities
- **Self-contained** - Single executable with embedded web interface, no external dependencies

## Quick Start

### Prerequisites

- Windows OS
- **Administrator privileges** (required for service management)

### Installation & Usage

1. Download `service-watch.exe` from releases

2. **Run as Administrator** (right-click → "Run as administrator"):

3. The application will:
   - Start quietly in your **system tray** (look for the Service Watch icon)
   - Begin monitoring on port 8080

4. **Access the dashboard**:
   - Right-click the system tray icon → "Open Dashboard"
   - Or browse to: `http://localhost:8080`

5. **View logs**:
   - Right-click the system tray icon → "Open Logs Folder"
   - Or check the `logs/` directory next to the executable

## System Tray Controls

Right-click the Service Watch system tray icon for quick access:

- **Open Dashboard** - Launch web interface in your browser
- **Open Logs Folder** - View log files in Windows Explorer  
- **Exit** - Close Service Watch completely

## Web Dashboard Features

### Service Management
- Add services to your watchlist with auto-discovery
- Start/stop/restart services with one click
- Configure auto-restart policies
- Real-time status monitoring with visual indicators

### Live Charts
- Pin services to charts for real-time CPU and memory monitoring
- Auto-scaling graphs that adapt to actual usage patterns
- Historical data visualization

### Event Logging
- Live log streaming with real-time filtering
- Historical log queries with search capabilities
- Multiple log levels and event types
- Export and analysis tools

## Configuration

### Port Configuration
The default port is `8080`. To change it:

1. Edit `main.go` and modify this line:
```go
addr := "127.0.0.1:8080"  // Change 8080 to your desired port
```

2. Rebuild the application:
```bash
go build -ldflags="-H windowsgui" -o service-watch.exe .
```

### Data Storage
- **Logs**: Stored in `logs/events.jsonl` (next to executable)
- **Configuration**: Stored in `watchlist.json` (next to executable)
- **Log Rotation**: Automatic (10MB max, 5 backups, 7 days retention)

### Auto-Start (Optional)
To start Service Watch automatically with Windows:

1. Copy `service-watch.exe` to a permanent location
2. Add to Windows startup folder: `Win+R` → `shell:startup`
3. Create a shortcut to the executable in the startup folder

## Development

### Building from Source

```bash
# Standard build (with console window)
go build -o service-watch.exe .

# GUI build (system tray only, recommended)
go build -ldflags="-H windowsgui" -o service-watch.exe .
```

### Dependencies
- [chi](https://github.com/go-chi/chi) - HTTP router
- [systray](https://github.com/getlantern/systray) - System tray integration
- [gopsutil](https://github.com/shirou/gopsutil) - System metrics
- [lumberjack](https://github.com/natefinch/lumberjack) - Log rotation

### Project Structure
```
service-watch/
├── main.go              # Application entry point
├── internal/            # Go backend modules
├── web/                 # Svelte web dashboard source
├── dist/                # Built web assets (embedded)
├── logs/                # Event logs (created at runtime)
├── watchlist.json       # Watchlist configuration (created at runtime)
└── icon.ico            # System tray icon
```

## Troubleshooting

### "Access Denied" Errors
- Ensure you're running as Administrator
- Some Windows services require elevated privileges to manage

### Dashboard Won't Load
- Check if port 8080 is already in use
- Verify Windows Firewall isn't blocking the application
- Try accessing `http://127.0.0.1:8080` directly

### System Tray Icon Missing
- Check if the application is running in Task Manager
- Look in the "hidden icons" area of your system tray
- Restart the application as Administrator

## Event Types

Service Watch monitors and logs the following events:

- `app_started` - Application startup
- `host_resources` - System CPU/memory metrics  
- `service_status` - Service state and resource usage
- `restart_attempt` - Auto-restart initiated
- `restart_success` - Service restarted successfully
- `restart_failed` - Service restart failed
- `service_failed` - Service exceeded restart limits

## Platform Support

- ✅ **Windows 10/11** - Fully supported
- 🚧 **Linux** - Planned for future release

## Contributing

Contributions welcome! Please open an issue or submit a pull request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Roadmap

- [ ] Discord/Slack webhook notifications
- [ ] Performance alerting thresholds  
- [ ] Linux/systemd support
- [ ] Remote monitoring capabilities
