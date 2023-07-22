package dbrepo

import (
	"github.com/nanmenkaimak/user-management/internal/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) CreateUser(newUser models.Users) (string, error) {
	var userID string
	password, err := hashPassword(newUser.Password)
	if err != nil {
		return "", err
	}
	err = m.DB.Get(&userID,
		`insert into users (username, email, password) 
    			values ($1, $2, $3)
    			returning id`,
		newUser.Username, newUser.Email, password)
	if err != nil {
		return "", errors.Wrap(err, "insert")
	}

	return userID, nil
}

func (m *postgresDBRepo) Authenticate(email string, password string) (string, error) {
	var user models.LoginUser

	err := m.DB.Get(&user,
		`select id, password from users where email = $1`, email)
	if err != nil {
		return "", errors.Wrap(err, "select auth")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return "", errors.Wrap(err, "incorrect password")
	} else if err != nil {
		return "", errors.Wrap(err, "password auth")
	}

	return user.ID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
