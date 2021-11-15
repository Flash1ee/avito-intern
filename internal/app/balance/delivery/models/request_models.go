package models

type RequestTransfer struct {
	FromUserID int64 `json:"from"`
	ToUserID   int64 `json:"to"`
	Amount     int64 `json:"amount"`
}

type RequestBalance struct {
	UserID int64 `json:"id"`
}

type RequestUpdateBalance struct {
	Type   int64 `json:"operation"`
	Amount int64 `json:"amount"`
}
