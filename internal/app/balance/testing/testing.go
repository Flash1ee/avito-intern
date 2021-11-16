package testing

import (
	request_response_models "avito-intern/internal/app/balance/delivery/models"
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
func TestUpdateBalance(t *testing.T) *request_response_models.RequestUpdateBalance {
	t.Helper()
	return &request_response_models.RequestUpdateBalance{
		Type:   models.ADD_BALANCE,
		Amount: 100.2,
	}
}
