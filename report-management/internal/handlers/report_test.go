package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	mock_repository "github.com/nanmenkaimak/report-management/internal/repository/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestReportByDate(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockDatabaseRepo, startDate time.Time, endDate time.Time, userID string)

	type reportStr = struct {
		StartDate time.Time
		EndDate   time.Time
		UserID    string
	}

	testReport := []struct {
		name               string
		inputUser          reportStr
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedStatusBody string
	}{
		{
			name: "OK",
			inputUser: reportStr{
				StartDate: time.Date(2023, 07, 01, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 8, 01, 0, 0, 0, 0, time.UTC),
				UserID:    "65ea156d-026d-41f6-8f48-075dca910277",
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, startDate time.Time, endDate time.Time, userID string) {
				s.EXPECT().ReportByDate(startDate, endDate, userID).Return(50_000, 50_000, 100_000, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedStatusBody: `{"Total":50000,"Expense":50000,"Income":100000}`,
		},
		{
			name: "No end date",
			inputUser: reportStr{
				StartDate: time.Date(2023, 07, 01, 0, 0, 0, 0, time.UTC),
				UserID:    "65ea156d-026d-41f6-8f48-075dca910277",
			},
			mockBehavior:       func(s *mock_repository.MockDatabaseRepo, startDate time.Time, endDate time.Time, userID string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedStatusBody: `{"message":"parse end date"}`,
		},
		{
			name: "No start date",
			inputUser: reportStr{
				EndDate: time.Date(2023, 8, 01, 0, 0, 0, 0, time.UTC),
				UserID:  "65ea156d-026d-41f6-8f48-075dca910277",
			},
			mockBehavior:       func(s *mock_repository.MockDatabaseRepo, startDate time.Time, endDate time.Time, userID string) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedStatusBody: `{"message":"parse start date"}`,
		},
		{
			name: "Wrong uuid",
			inputUser: reportStr{
				StartDate: time.Date(2023, 7, 01, 0, 0, 0, 0, time.UTC),
				EndDate:   time.Date(2023, 8, 01, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(s *mock_repository.MockDatabaseRepo, startDate time.Time, endDate time.Time, userID string) {
				s.EXPECT().ReportByDate(startDate, endDate, userID).Return(0, 0, 0, errors.Wrap(errors.New("some error"), "select expense amount month"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedStatusBody: `{"message":"select expense amount month: some error"}`,
		},
	}

	for _, tt := range testReport {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			newRepo := mock_repository.NewMockDatabaseRepo(c)
			tt.mockBehavior(newRepo, tt.inputUser.StartDate, tt.inputUser.EndDate, tt.inputUser.UserID)

			repo := &Repository{newRepo}
			NewHandlers(repo)

			r := chi.NewRouter()
			r.Get("/nmk/report", Repo.ReportByDate)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/nmk/report", nil)
			q := req.URL.Query()
			if !tt.inputUser.StartDate.IsZero() {
				q.Add("start", tt.inputUser.StartDate.Format("2006-01-02"))
			}
			if !tt.inputUser.EndDate.IsZero() {
				q.Add("end", tt.inputUser.EndDate.Format("2006-01-02"))
			}
			req.URL.RawQuery = q.Encode()

			ctx := context.WithValue(req.Context(), "id", tt.inputUser.UserID)
			req = req.WithContext(ctx)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedStatusBody, w.Body.String())
		})
	}
}
