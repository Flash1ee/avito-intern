package models

const (
	DIFF_BALANCE = iota
	ADD_BALANCE
)

type Balance struct {
	ID     int64
	Amount int64
}
type TransferMoney struct {
	SenderID        int64
	SenderBalance   int64
	ReceiverID      int64
	ReceiverBalance int64
}
