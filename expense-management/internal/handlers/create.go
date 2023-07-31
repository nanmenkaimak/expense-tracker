package handlers

import (
	"encoding/json"
	"github.com/nanmenkaimak/expense-management/internal/models"
	"github.com/nanmenkaimak/expense-management/internal/rabbit"
	"net/http"
)

func (m *Repository) CreateExpenses(w http.ResponseWriter, r *http.Request) {
	var newExpense models.Expenses

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&newExpense); err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("id").(string)
	newExpense.UserID = userID

	id, err := m.DB.CreateExpense(newExpense)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	type responseCreate struct {
		ID string `json:"id"`
	}

	renderJSON(w, responseCreate{
		ID: id,
	})
}

type message struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (m *Repository) Expenses(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	data, err := rabbit.ReceiveMessage(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var messageRes message
	json.Unmarshal([]byte(data), &messageRes)
	renderJSON(w, messageRes)
}
