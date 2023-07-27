package dbrepo

import (
	"github.com/pkg/errors"
	"time"
)

func (m *postgresDBRepo) ReportByDate(start time.Time, end time.Time, userID string) (int, int, int, error) {
	var expensesAmount int
	var incomeAmount int
	err := m.DB.Get(&expensesAmount,
		`select sum(expenses.amount) from expenses inner join categories on expenses.category_id = categories.id
				where categories.is_income = $1 
				  and expenses.date_time between $2 and $3 and user_id = $4`, false, start, end, userID)
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "select expense amount month")
	}
	err = m.DB.Get(&incomeAmount,
		`select sum(expenses.amount) from expenses inner join categories on expenses.category_id = categories.id
				where categories.is_income = $1
				  and expenses.date_time between $2 and $3 and user_id = $4`, true, start, end, userID)
	if err != nil {
		return 0, 0, 0, errors.Wrap(err, "select income amount month")
	}

	return absDiffInt(expensesAmount, incomeAmount), expensesAmount, incomeAmount, nil
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return -1 * (x - y)
}
