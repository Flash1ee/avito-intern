package transaction

import "avito-intern/internal/app/transaction/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mock_transaction -mock_names=Usecase=TransactionUsecase . Usecase

type Usecase interface {
	// GetTransactions Errors:
	//		transaction_repository.NotFound
	//		app.GeneralError with Errors:
	//			transaction_repository.DefaultErrDB
	GetTransactions(userID int64, paginator *models.Paginator) ([]models.Transaction, error)
}
