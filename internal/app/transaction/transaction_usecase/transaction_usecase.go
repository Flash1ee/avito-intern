package transaction_usecase

import (
	"avito-intern/internal/app/transaction"
	"avito-intern/internal/app/transaction/models"
)

type TransactionUsecase struct {
	repo transaction.Repository
}

func NewTransactionUsecase(repo transaction.Repository) *TransactionUsecase {
	return &TransactionUsecase{
		repo: repo,
	}
}

// GetTransactions Errors:
//		transaction_repository.NotFound
//		app.GeneralError with Errors:
//			transaction_repository.DefaultErrDB
func (uc *TransactionUsecase) GetTransactions(userID int64, paginator *models.Paginator) ([]models.Transaction, error) {
	return uc.repo.GetTransactions(userID, paginator)
}
