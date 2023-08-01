package repository

import "time"

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type DatabaseRepo interface {
	ReportByDate(start time.Time, end time.Time, userID string) (int, int, int, error)
}
