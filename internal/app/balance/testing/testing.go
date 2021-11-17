package testing

import (
	request_response_models "avito-intern/internal/app/balance/delivery/models"
	"avito-intern/internal/app/balance/models"
	transaction_constants "avito-intern/internal/app/models"
	transaction_models "avito-intern/internal/app/transaction/models"
	"fmt"
	"testing"
	"time"
)

func TestBalance(t *testing.T) *models.Balance {
	t.Helper()
	return &models.Balance{
		ID:     1,
		Amount: 10.2,
	}
}
func TestTransaction(t *testing.T) *models.TransferMoney {
	t.Helper()
	return &models.TransferMoney{
		SenderID:        1,
		SenderBalance:   10.2,
		ReceiverID:      2,
		ReceiverBalance: 5,
	}
}
func TestUpdateBalance(t *testing.T) *request_response_models.RequestUpdateBalance {
	t.Helper()
	return &request_response_models.RequestUpdateBalance{
		Type:   models.ADD_BALANCE,
		Amount: 100.2,
	}
}
func TestAddBalanceDescription(id int64, operation string, t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("user <id = %d> %s account", id, operation)
}
func TestTransferDescription(senderID int64, receiverID int64, t *testing.T) string {
	t.Helper()
	return fmt.Sprintf("user <id = %d> send money to user <id = %d>", senderID, receiverID)
}
func TestPaginator(t *testing.T) *transaction_models.Paginator {
	t.Helper()
	return &transaction_models.Paginator{
		Page:          1,
		Count:         5,
		SortDirection: transaction_models.NO_DIRECTION,
		SortField:     transaction_models.NO_ORDER,
	}
}
func TestTransactions(t *testing.T) []transaction_models.Transaction {
	t.Helper()
	return []transaction_models.Transaction{
		{
			UserID:      1,
			ReceiverID:  2,
			Type:        transaction_constants.TextTransactionToType[transaction_constants.WRITE_OFF],
			CreatedAt:   time.Now(),
			Description: "decription of write-off",
		},
		{
			UserID:      2,
			ReceiverID:  1,
			Type:        transaction_constants.TextTransactionToType[transaction_constants.TRANSFER],
			CreatedAt:   time.Now(),
			Description: "decription of transfer",
		},
	}
}
