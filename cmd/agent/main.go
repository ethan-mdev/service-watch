package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ethan-mdev/service-watch/internal/handlers"
	"github.com/ethan-mdev/service-watch/internal/platform"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// v1 routes
	r.Route("/v1", func(v1 chi.Router) {
		v1.Get("/status", handlers.Status())
		svcHTTP := handlers.NewServiceHTTP(platform.MakeServiceManager())
		v1.Mount("/services", svcHTTP.Routes())
	})

	srv := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("service-watch listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
