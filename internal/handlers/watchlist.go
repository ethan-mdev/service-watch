package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ethan-mdev/service-watch/internal/core"
	"github.com/ethan-mdev/service-watch/internal/utils"
	"github.com/go-chi/chi/v5"
)

type WatchlistHTTP struct {
	M core.WatchlistManager
}

func NewWatchlistHTTP(m core.WatchlistManager) *WatchlistHTTP {
	return &WatchlistHTTP{M: m}
}

// Routes sets up the HTTP routes for watchlist management.
func (h *WatchlistHTTP) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Post("/", h.add)
	r.Route("/{name}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Put("/", h.update)
		r.Delete("/", h.remove)
	})
	return r
}

func (h *WatchlistHTTP) list(w http.ResponseWriter, r *http.Request) {
	items, err := h.M.List(r.Context())
	if err != nil {
		utils.RespondWithError(w, 500, "failed to list watchlist", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"items": items})
}

func (h *WatchlistHTTP) get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	item, err := h.M.Get(r.Context(), name)
	if err != nil {
		utils.RespondWithError(w, 404, "watchlist item not found", err)
		return
	}
	utils.RespondWithJSON(w, 200, item)
}

func (h *WatchlistHTTP) add(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ServiceName string `json:"serviceName"`
		AutoRestart bool   `json:"autoRestart"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, 400, "invalid request body", err)
		return
	}

	if req.ServiceName == "" {
		utils.RespondWithError(w, 400, "serviceName is required", nil)
		return
	}

	if err := h.M.Add(r.Context(), req.ServiceName, req.AutoRestart); err != nil {
		utils.RespondWithError(w, 400, "failed to add to watchlist", err)
		return
	}
	utils.RespondWithJSON(w, 201, map[string]any{"added": true})
}

func (h *WatchlistHTTP) update(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var req struct {
		AutoRestart bool `json:"autoRestart"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, 400, "invalid request body", err)
		return
	}

	if err := h.M.Update(r.Context(), name, req.AutoRestart); err != nil {
		utils.RespondWithError(w, 400, "failed to update watchlist item", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"updated": true})
}

func (h *WatchlistHTTP) remove(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := h.M.Remove(r.Context(), name); err != nil {
		utils.RespondWithError(w, 404, "failed to remove from watchlist", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"removed": true})
}
