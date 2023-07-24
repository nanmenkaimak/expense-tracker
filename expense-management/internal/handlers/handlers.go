package handlers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/nanmenkaimak/expense-management/internal/repository"
	"github.com/nanmenkaimak/expense-management/internal/repository/dbrepo"
	"net/http"
)

// Repo is repository used by handlers
var Repo *Repository

// Repository is repository type
type Repository struct {
	DB repository.DatabaseRepo
}

// NewRepo creates new repository
func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		DB: dbrepo.NewPostgresRepo(db),
	}
}

// NewHandlers sets repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) []byte {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return js
}
