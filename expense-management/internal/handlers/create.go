package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/nanmenkaimak/expense-management/internal/models"
	"github.com/nanmenkaimak/expense-management/internal/rabbit"
	"net/http"
)

func (m *Repository) CreateExpenses(w http.ResponseWriter, r *http.Request) {
	var newExpense models.Expenses

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&newExpense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := m.DB.CreateExpense(newExpense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, newExpense)
}

type message struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	data, err := rabbit.ReceiveMessage(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var messageRes message
	json.Unmarshal([]byte(data), &messageRes)
	fmt.Println(messageRes.ID)
}
