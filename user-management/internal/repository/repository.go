package repository

import "github.com/nanmenkaimak/user-management/internal/models"

type DatabaseRepo interface {
	CreateUser(newUser models.Users) (string, error)
	Authenticate(email string, password string) (string, error)
}
