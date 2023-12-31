package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/nanmenkaimak/user-management/internal/models"
	"github.com/pkg/errors"
	"net/http"
	"unicode"
)

func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.Users

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&newUser); err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	ok := govalidator.IsEmail(newUser.Email)
	if !ok {
		newErrorResponse(w, errorResponse{Message: errors.New("invalid email").Error()}, http.StatusBadRequest)
		return
	}
	err := validPassword(newUser.Password)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusBadRequest)
		return
	}

	id, err := m.DB.CreateUser(newUser)
	if err != nil {
		newErrorResponse(w, errorResponse{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	type responseSignUp struct {
		ID string `json:"id"`
	}

	renderJSON(w, responseSignUp{
		ID: id,
	})
}

func validPassword(s string) error {
next:
	for name, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return fmt.Errorf("password must have at least one %s character", name)
	}
	return nil
}
