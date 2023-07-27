package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nanmenkaimak/report-management/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/home/{username}", handlers.Repo.Home)

	mux.Route("/{username}", func(r chi.Router) {
		r.Use(Auth)
		r.Get("/report", handlers.Repo.ReportByDate)
	})

	return mux
}
