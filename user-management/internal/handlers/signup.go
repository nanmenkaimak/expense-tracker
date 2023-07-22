package handlers

import (
	"encoding/json"
	"github.com/nanmenkaimak/user-management/internal/models"
	"net/http"
)

func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.Users

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := m.DB.CreateUser(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, newUser)
}
