package handlers

import (
	"net/http"
	"time"

	"github.com/ethan-mdev/service-watch/internal/system"
	"github.com/ethan-mdev/service-watch/internal/utils"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type StatusHandler struct {
	sys *system.System
}

func NewStatusHandler(sys *system.System) *StatusHandler { return &StatusHandler{sys: sys} }

// GET /v1/status
func (s *StatusHandler) Health(w http.ResponseWriter, r *http.Request) {
	cpuPercent, _ := cpu.Percent(0, false) // snapshot
	memInfo, _ := mem.VirtualMemory()

	response := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service": map[string]any{
			"uptime": time.Since(s.sys.Start).String(),
		},
		"host": map[string]any{
			"hostname":       s.sys.HostInfo.Hostname,
			"os":             s.sys.HostInfo.OS,             // "linux" | "windows" | "darwin"
			"platform":       s.sys.HostInfo.Platform,       // "ubuntu", "windows", etc.
			"platformFamily": s.sys.HostInfo.PlatformFamily, // "debian", "windows", etc.
			"platformVer":    s.sys.HostInfo.PlatformVersion,
			"kernel":         s.sys.HostInfo.KernelVersion,
			"systemUptime":   s.sys.HostInfo.Uptime, // seconds
		},
		"cpu": map[string]any{
			"cores":   len(cpuPercent),
			"percent": cpuPercent,
		},
		"memory": map[string]any{
			"total":     memInfo.Total,
			"available": memInfo.Available,
			"used":      memInfo.Used,
			"usedPct":   memInfo.UsedPercent,
		},
	}
	utils.RespondWithJSON(w, http.StatusOK, response)
}
