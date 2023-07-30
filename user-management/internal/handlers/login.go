package handlers

import (
	"encoding/json"
	"github.com/nanmenkaimak/user-management/internal/JWT"
	"github.com/nanmenkaimak/user-management/internal/models"
	"github.com/nanmenkaimak/user-management/internal/rabbit"
	"net/http"
)

type message struct {
	ID    string
	Token string
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	var user models.Users

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&user); err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	id, username, err := m.DB.Authenticate(user.Email, user.Password)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	token, err := JWT.GenerateToken(id)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	messageSend := message{ID: id, Token: token}

	messageJSON := renderJSON(w, messageSend)
	err = rabbit.SendMessage(messageJSON, username)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
}
