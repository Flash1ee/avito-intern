package models

import "time"

const (
	NO_ORDER = iota
	DATE
	SUM
	NO_DIRECTION = iota + 10
	ASC
	DESC
)

var TransactionQueryParams = map[int]string{
	NO_ORDER:     "",
	DATE:         "created_at",
	SUM:          "amount",
	NO_DIRECTION: "",
	ASC:          "asc",
	DESC:         "desc",
}

type Transaction struct {
	UserID      int64     `json:"sender_id"`
	ReceiverID  int64     `json:"receiver_id,omitempty"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" swaggertype:"integer"`
	Amount      float64   `json:"amount"`
}

type Paginator struct {
	Page          int `json:"page"`
	Count         int `json:"count"`
	SortField     int `json:"sort_field"`
	SortDirection int `json:"sort_direction"`
}
