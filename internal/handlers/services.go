package handlers

import (
	"net/http"

	"github.com/ethan-mdev/service-watch/internal/core"
	"github.com/ethan-mdev/service-watch/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ServiceHTTP struct{ M core.ServiceManager }

func NewServiceHTTP(m core.ServiceManager) *ServiceHTTP { return &ServiceHTTP{M: m} }

// Routes sets up the HTTP routes for service management.
func (h *ServiceHTTP) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", h.list)
	r.Route("/{name}", func(r chi.Router) {
		r.Get("/", h.get)
		r.Post("/start", h.start)
		r.Post("/stop", h.stop)
		r.Post("/restart", h.restart)
	})
	return r
}

func (h *ServiceHTTP) list(w http.ResponseWriter, r *http.Request) {
	rows, err := h.M.List(r.Context())
	if err != nil {
		utils.RespondWithError(w, 500, "failed to list services", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"items": rows})
}

func (h *ServiceHTTP) get(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	row, err := h.M.Get(r.Context(), name)
	if err != nil {
		utils.RespondWithError(w, 500, "failed to query service", err)
		return
	}
	utils.RespondWithJSON(w, 200, row)
}

func (h *ServiceHTTP) start(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := h.M.Start(r.Context(), name); err != nil {
		utils.RespondWithError(w, 424, "start failed", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"accepted": true})
}

func (h *ServiceHTTP) stop(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := h.M.Stop(r.Context(), name); err != nil {
		utils.RespondWithError(w, 424, "stop failed", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"accepted": true})
}

func (h *ServiceHTTP) restart(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if err := h.M.Restart(r.Context(), name); err != nil {
		utils.RespondWithError(w, 424, "restart failed", err)
		return
	}
	utils.RespondWithJSON(w, 200, map[string]any{"accepted": true})
}
