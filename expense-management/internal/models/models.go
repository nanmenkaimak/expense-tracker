package models

import "time"

type Expenses struct {
	ID          string    `json:"id"`
	Amount      int       `json:"amount"`
	CategoryID  int       `json:"category_id"`
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
