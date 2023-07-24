package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nanmenkaimak/expense-management/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Route("/{username}", func(r chi.Router) {
		r.Get("/", handlers.Repo.Expenses)
		r.Post("/new", handlers.Repo.CreateExpenses)
		r.Delete("/delete", handlers.Repo.DeleteExpense)
	})

	return mux
}
