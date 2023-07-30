package repository

import "github.com/nanmenkaimak/user-management/internal/models"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type DatabaseRepo interface {
	CreateUser(newUser models.Users) (string, error)
	Authenticate(email string, password string) (string, string, error)
}
