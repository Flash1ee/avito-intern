package models

type ErrResponse struct {
	Err string `json:"error"`
}

type ResponseBalance struct {
	UserID  int64   `json:"user_id,omitempty"`
	Balance float64 `json:"balance"`
}
type ResponseTransfer struct {
	SenderID        int64   `json:"sender_id"`
	SenderBalance   float64 `json:"sender_balance"`
	ReceiverID      int64   `json:"receiver_id"`
	ReceiverBalance float64 `json:"receiver_balance"`
}
