package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/nanmenkaimak/expense-management/internal/models"
	mock_repository "github.com/nanmenkaimak/expense-management/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateExpenses(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockDatabaseRepo, expense models.Expenses)

	testCreate := []struct {
		name               string
		inputBody          string
		inputUser          models.Expenses
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedStatusBody string
	}{
		{
			name:      "OK",
			inputBody: `{"amount":50000, "category_id":2}`,
			inputUser: models.Expenses{
				UserID:     "65ea156d-026d-41f6-8f48-075dca910277",
				Amount:     50000,
				CategoryID: 2,
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, expense models.Expenses) {
				s.EXPECT().CreateExpense(expense).Return("56cb90b6-7432-4eb7-bce9-62ed4f179c60", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedStatusBody: `{"id":"56cb90b6-7432-4eb7-bce9-62ed4f179c60"}`,
		},
		{
			name:               "wrong input json",
			inputBody:          `{"apple":"madikhanmeshok"}`,
			mockBehavior:       func(s *mock_repository.MockDatabaseRepo, expense models.Expenses) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedStatusBody: `{"message":"json: unknown field \"apple\""}`,
		},
		{
			name:      "wrong category id",
			inputBody: `{"amount":50000, "category_id":3}`,
			inputUser: models.Expenses{
				UserID:     "65ea156d-026d-41f6-8f48-075dca910277",
				Amount:     50000,
				CategoryID: 3,
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, expense models.Expenses) {
				s.EXPECT().CreateExpense(expense).Return("", errors.Wrap(errors.New("pq: INSERT или UPDATE в таблице \"expenses\" нарушает ограничение внешнего ключа \"expenses_category_id_fkey\""), "insert expense"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedStatusBody: `{"message":"insert expense: pq: INSERT или UPDATE в таблице \"expenses\" нарушает ограничение внешнего ключа \"expenses_category_id_fkey\""}`,
		},
	}

	for _, tt := range testCreate {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newRepo := mock_repository.NewMockDatabaseRepo(c)
			tt.mockBehavior(newRepo, tt.inputUser)

			repo := &Repository{newRepo}
			NewHandlers(repo)

			r := chi.NewRouter()
			r.Post("/nmk/new", Repo.CreateExpenses)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/nmk/new",
				bytes.NewBufferString(tt.inputBody))

			ctx := context.WithValue(req.Context(), "id", tt.inputUser.UserID)
			req = req.WithContext(ctx)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedStatusBody, w.Body.String())
		})
	}
}
