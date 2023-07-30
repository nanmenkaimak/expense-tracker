package handlers

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"github.com/nanmenkaimak/user-management/internal/models"
	mock_repository "github.com/nanmenkaimak/user-management/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockDatabaseRepo, user models.Users)

	testSignUp := []struct {
		name                string
		inputBody           string
		inputUser           models.Users
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"nanmenkaimak", "email":"a@gmail.com", "password":"Qwertyuiop1*"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "a@gmail.com",
				Password: "Qwertyuiop1*",
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, user models.Users) {
				s.EXPECT().CreateUser(user).Return("57858cc4-d28c-47ea-a2f8-9dac35bd4d6e", nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: `{"id":"57858cc4-d28c-47ea-a2f8-9dac35bd4d6e"}`,
		},
		{
			name:      "invalid email",
			inputBody: `{"username":"nanmenkaimak", "email":"einebashnabashbarma", "password":"Qwertyuiop1*"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "einebashnabashbarma",
				Password: "Qwertyuiop1*",
			},
			mockBehavior:        func(s *mock_repository.MockDatabaseRepo, user models.Users) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"invalid email"}`,
		},
		{
			name:      "uppercase",
			inputBody: `{"username":"nanmenkaimak", "email":"a@gmail.com", "password":"qwertyuiop1*"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "a@gmail.com",
				Password: "qwertyuiop1*",
			},
			mockBehavior:        func(s *mock_repository.MockDatabaseRepo, user models.Users) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"password must have at least one upper case character"}`,
		},
		{
			name:      "lowercase",
			inputBody: `{"username":"nanmenkaimak", "email":"a@gmail.com", "password":"QWERTYUIOP1*"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "a@gmail.com",
				Password: "QWERTYUIOP1*",
			},
			mockBehavior:        func(s *mock_repository.MockDatabaseRepo, user models.Users) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"password must have at least one lower case character"}`,
		},
		{
			name:      "numeric",
			inputBody: `{"username":"nanmenkaimak", "email":"a@gmail.com", "password":"Qwertyuiop*"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "a@gmail.com",
				Password: "Qwertyuiop*",
			},
			mockBehavior:        func(s *mock_repository.MockDatabaseRepo, user models.Users) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"password must have at least one numeric character"}`,
		},
		{
			name:      "special",
			inputBody: `{"username":"nanmenkaimak", "email":"a@gmail.com", "password":"Qwertyuiop1"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "a@gmail.com",
				Password: "Qwertyuiop1",
			},
			mockBehavior:        func(s *mock_repository.MockDatabaseRepo, user models.Users) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"password must have at least one special character"}`,
		},
		{
			name:      "similar username",
			inputBody: `{"username":"nanmenkaimak", "email":"a@gmail.com", "password":"Qwertyuiop1*"}`,
			inputUser: models.Users{
				Username: "nanmenkaimak",
				Email:    "a@gmail.com",
				Password: "Qwertyuiop1*",
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, user models.Users) {
				s.EXPECT().CreateUser(user).Return("", errors.New("insert: pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_username_key\""))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"insert: pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_username_key\""}`,
		},
		{
			name:      "similar email",
			inputBody: `{"username":"nanmenkaima", "email":"a@gmail.com", "password":"Qwertyuiop1*"}`,
			inputUser: models.Users{
				Username: "nanmenkaima",
				Email:    "a@gmail.com",
				Password: "Qwertyuiop1*",
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, user models.Users) {
				s.EXPECT().CreateUser(user).Return("", errors.New("insert: pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_email_key\""))
			},
			expectedStatusCode:  http.StatusInternalServerError,
			expectedRequestBody: `{"message":"insert: pq: повторяющееся значение ключа нарушает ограничение уникальности \"users_email_key\""}`,
		},
		{
			name:                "wrong input json",
			inputBody:           `{"apple":"madikhanmeshok"}`,
			mockBehavior:        func(s *mock_repository.MockDatabaseRepo, user models.Users) {},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: `{"message":"json: unknown field \"apple\""}`,
		},
	}

	for _, testCase := range testSignUp {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			signup := mock_repository.NewMockDatabaseRepo(c)
			testCase.mockBehavior(signup, testCase.inputUser)

			repo := &Repository{signup}
			NewHandlers(repo)
			r := chi.NewRouter()
			r.Post("/auth/signup", Repo.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/signup",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}

}
