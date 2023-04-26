package domain

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type OrderStatus string

const (
	OrderNew          OrderStatus = "NEW"
	OrderInProcessing OrderStatus = "PROCESSING"
	OrderInvalid      OrderStatus = "INVALID"
	OrderProcessed    OrderStatus = "PROCESSED"
)

type Order struct {
	Number     string      `json:"number"`
	UserID     string      `json:"-" db:"user_id"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual"`
	UploadedAt UploadedAt  `json:"uploaded_at" db:"uploaded_at"`
}

type UploadedAt struct {
	time.Time
}

func (c *UploadedAt) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(time.RFC3339))), nil
}

func (c *UploadedAt) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return c.Time, nil
}

func (c *UploadedAt) Scan(value any) error {
	if value == nil {
		c.Time = time.Time{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("invalid UploadedAt value")
	}
	c.Time = t
	return nil
}
