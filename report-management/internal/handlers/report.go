package handlers

import (
	"github.com/pkg/errors"
	"net/http"
	"time"
)

func (m *Repository) ReportByDate(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("id").(string)
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	layout := "2006-01-02"

	startTime, err := time.Parse(layout, startDate)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: errors.New("parse start date").Error()}, http.StatusBadRequest)
		return
	}

	endTime, err := time.Parse(layout, endDate)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: errors.New("parse end date").Error()}, http.StatusBadRequest)
		return
	}

	total, expense, income, err := m.DB.ReportByDate(startTime, endTime, userID)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	type reportMoney struct {
		Total   int
		Expense int
		Income  int
	}

	renderJSON(w, reportMoney{
		Total:   total,
		Expense: expense,
		Income:  income,
	})
}
