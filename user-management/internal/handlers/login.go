package handlers

import (
	"encoding/json"
	"github.com/nanmenkaimak/user-management/internal/JWT"
	"github.com/nanmenkaimak/user-management/internal/models"
	"github.com/nanmenkaimak/user-management/internal/rabbit"
	"net/http"
)

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	var user models.LoginUser

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := m.DB.Authenticate(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := JWT.GenerateToken(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renderJSON(w, token)
	err = rabbit.SendMessage([]byte(token), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
