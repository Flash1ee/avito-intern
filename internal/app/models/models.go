package transaction_models

const (
	WRITE_OFF = iota
	REFILL
	TRANSFER
)

var TransactionTypeToText = map[string]int{
	"write-off": WRITE_OFF,
	"refill":    REFILL,
	"transfer":  TRANSFER,
}
var TextTransactionToType = map[int]string{
	WRITE_OFF: "write-off",
	REFILL:    "refill",
	TRANSFER:  "transfer",
}
