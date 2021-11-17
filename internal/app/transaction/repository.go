package transaction

import "avito-intern/internal/app/transaction/models"

//go:generate mockgen -destination=mocks/repository.go -package=mock_transaction -mock_names=Repository=TransactionRepository . Repository

type Repository interface {
	// GetTransactions Errors:
	//		transaction_repository.NotFound
	//		app.GeneralError with Errors:
	//			transaction_repository.DefaultErrDB
	GetTransactions(userID int64, paginator *models.Paginator) ([]models.Transaction, error)
}
