package repository

import "time"

type DatabaseRepo interface {
	ReportByDate(start time.Time, end time.Time, userID string) (int, int, int, error)
}
