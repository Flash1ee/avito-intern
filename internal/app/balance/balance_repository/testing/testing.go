package testing

import (
	"avito-intern/internal/app/balance/models"
	"testing"
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
