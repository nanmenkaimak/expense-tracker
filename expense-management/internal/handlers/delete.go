package handlers

import (
	"github.com/pkg/errors"
	"net/http"
)

func (m *Repository) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		http.Error(w, errors.New("cannot get id of expense").Error(), http.StatusBadRequest)
		return
	}

	ok, err := m.DB.DeleteExpense(expenseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !ok {
		http.Error(w, errors.New("no rows is affected").Error(), http.StatusBadRequest)
		return
	}

}
