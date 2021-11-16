package models

const (
	DIFF_BALANCE = iota
	ADD_BALANCE
)

type Balance struct {
	ID     int64
	Amount float64
}
type TransferMoney struct {
	SenderID        int64
	SenderBalance   float64
	ReceiverID      int64
	ReceiverBalance float64
}
