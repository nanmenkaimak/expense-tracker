package dbrepo

import (
	"github.com/nanmenkaimak/expense-management/internal/models"
	"github.com/pkg/errors"
	"time"
)

func (m *postgresDBRepo) CreateExpense(newExpense models.Expenses) (string, error) {
	var expenseID string
	if newExpense.Date.IsZero() {
		newExpense.Date = time.Now()
	}
	err := m.DB.Get(&expenseID,
		`insert into expenses (amount, category_id, user_id, description, date_time)
				values ($1, $2, $3, $4, $5) returning id`,
		newExpense.Amount, newExpense.CategoryID, newExpense.UserID, newExpense.Description, newExpense.Date)
	if err != nil {
		return "", errors.Wrap(err, "insert expense")
	}

	return expenseID, nil
}

func (m *postgresDBRepo) DeleteExpense(expenseID string) (bool, error) {
	res, err := m.DB.Exec(`delete from expenses where id = $1`, expenseID)
	if err != nil {
		return false, errors.Wrap(err, "delete expense")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, errors.Wrap(err, "rows affected")
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
