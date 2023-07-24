package repository

import "github.com/nanmenkaimak/expense-management/internal/models"

type DatabaseRepo interface {
	CreateExpense(newExpense models.Expenses) (string, error)
	DeleteExpense(expenseID string) (bool, error)
}
