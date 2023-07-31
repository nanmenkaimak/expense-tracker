package repository

import "github.com/nanmenkaimak/expense-management/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type DatabaseRepo interface {
	CreateExpense(newExpense models.Expenses) (string, error)
	DeleteExpense(expenseID string) (bool, error)
}
