package handlers

import (
	"encoding/json"
	"github.com/nanmenkaimak/report-management/internal/rabbit"
	"github.com/pkg/errors"
	"net/http"
)

type message struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	data, err := rabbit.ReceiveMessage(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = rabbit.SendMessage([]byte(data), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var messageRes message
	err = json.Unmarshal([]byte(data), &messageRes)
	if err != nil {
		http.Error(w, errors.Wrap(err, "unmarshal json home").Error(), http.StatusBadRequest)
	}
	renderJSON(w, messageRes)
}
