package models

import "time"

type Transaction struct {
	UserID      int64     `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Amount      float64   `json:"amount"`
}

type Paginator struct {
	Page          int    `json:"page"`
	Count         int    `json:"count"`
	SortField     string `json:"sort_field"`
	SortDirection string `json:"sort_direction"`
}
