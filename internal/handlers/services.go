package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ethan-mdev/service-watch/internal/system"
)

type ServiceHandler struct {
	sys *system.System
}

func NewServiceHandler(sys *system.System) *ServiceHandler { return &ServiceHandler{sys: sys} }

// GET /v1/services
func (h *ServiceHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	items, err := h.sys.Services.ListServices(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"hostOS":   h.sys.HostInfo.OS,
		"count":    len(items),
		"services": items,
	})
}
