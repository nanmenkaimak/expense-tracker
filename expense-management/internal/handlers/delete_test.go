package handlers

import (
	"github.com/go-chi/chi/v5"
	mock_repository "github.com/nanmenkaimak/expense-management/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteExpense(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockDatabaseRepo, expenseID string)

	testDelete := []struct {
		name               string
		inputUser          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedStatusBody string
	}{
		{
			name:      "OK",
			inputUser: "4667e606-e333-4557-b502-8374abac38ce",
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, expenseID string) {
				s.EXPECT().DeleteExpense("4667e606-e333-4557-b502-8374abac38ce").Return(true, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedStatusBody: "",
		},
		{
			name:               "missing ID",
			mockBehavior:       func(s *mock_repository.MockDatabaseRepo, expenseID string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedStatusBody: `{"message":"cannot get id of expense"}`,
		},
		{
			name:      "wrong uuid",
			inputUser: "ad4dd7d2-5556-4e77-9245-9ead5429f52",
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, expenseID string) {
				s.EXPECT().DeleteExpense("ad4dd7d2-5556-4e77-9245-9ead5429f52").Return(false, errors.Wrap(errors.New(`pq: неверный синтаксис для типа uuid: "ad4dd7d2-5556-4e77-9245-9ead5429f52"`), "delete expense"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedStatusBody: `{"message":"delete expense: pq: неверный синтаксис для типа uuid: \"ad4dd7d2-5556-4e77-9245-9ead5429f52\""}`,
		},
		{
			name:      "no rows affected",
			inputUser: "4667e606-e333-4557-b502-8374abac38ce",
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, expenseID string) {
				s.EXPECT().DeleteExpense("4667e606-e333-4557-b502-8374abac38ce").Return(false, nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedStatusBody: `{"message":"no rows is affected"}`,
		},
	}
	for _, tt := range testDelete {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newRepo := mock_repository.NewMockDatabaseRepo(c)
			tt.mockBehavior(newRepo, tt.inputUser)

			repo := &Repository{newRepo}
			NewHandlers(repo)

			r := chi.NewRouter()
			r.Delete("/nmk/delete", Repo.DeleteExpense)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/nmk/delete", nil)
			q := req.URL.Query()
			if len(tt.inputUser) > 0 {
				q.Add("id", tt.inputUser)
			}
			req.URL.RawQuery = q.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedStatusBody, w.Body.String())
		})
	}
}
