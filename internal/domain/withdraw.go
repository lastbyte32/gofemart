package domain

import "time"

type Withdraw struct {
	ID          int64     `json:"-" db:"id"`
	UserID      string    `json:"-" db:"user_id"`
	OrderNumber string    `json:"order" db:"order_number"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}
