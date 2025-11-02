package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ethan-mdev/service-watch/internal/handlers"
	"github.com/ethan-mdev/service-watch/internal/system"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	sys, err := system.InitSystem()
	if err != nil {
		log.Fatalf("init system: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	status := handlers.NewStatusHandler(sys)
	svcs := handlers.NewServiceHandler(sys)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/status", status.Health)
		r.Get("/services", svcs.ListServices)
	})

	port := ":8080"
	fmt.Printf("Server starting on %s (OS=%s, Platform=%s)\n",
		port, sys.HostInfo.OS, sys.HostInfo.Platform)
	log.Fatal(http.ListenAndServe(port, r))
}
