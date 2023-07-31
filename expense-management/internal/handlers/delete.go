package handlers

import (
	"github.com/pkg/errors"
	"net/http"
)

func (m *Repository) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	expenseID := r.URL.Query().Get("id")
	if expenseID == "" {
		newErrorResponse(w, errorResponse{Message: errors.New("cannot get id of expense").Error()}, http.StatusBadRequest)
		return
	}

	ok, err := m.DB.DeleteExpense(expenseID)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	if !ok {
		newErrorResponse(w, errorResponse{Message: errors.New("no rows is affected").Error()}, http.StatusBadRequest)
		return
	}

}
