package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nanmenkaimak/user-management/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Route("/auth", func(r chi.Router) {
		r.Post("/signup", handlers.Repo.SignUp)
		r.Post("/login", handlers.Repo.Login)
	})

	return mux
}
