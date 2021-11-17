package models

import (
	transaction_models "avito-intern/internal/app/models"
	"avito-intern/internal/app/transaction/models"
)

type ErrResponse struct {
	Err string `json:"error"`
}

type NotFoundResponse struct {
	Ok string `json:"OK"`
}
type ResponseTransactions struct {
	Transactions []models.Transaction `json:"transactions"`
}

func ToResponseTransactions(transactions []models.Transaction) ResponseTransactions {
	res := make([]models.Transaction, 0, len(transactions))

	for _, t := range transactions {
		response := models.Transaction{
			Type:        t.Type,
			UserID:      t.UserID,
			CreatedAt:   t.CreatedAt,
			Amount:      t.Amount,
			Description: t.Description,
		}
		if transaction_models.TransactionTypeToText[t.Type] == transaction_models.TRANSFER {
			response.ReceiverID = t.ReceiverID
		}
		res = append(res, response)
	}
	return ResponseTransactions{
		Transactions: res,
	}
}
