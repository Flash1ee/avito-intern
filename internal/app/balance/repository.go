package balance

import "avito-intern/internal/app/balance/models"

//go:generate mockgen -destination=mocks/repository.go -package=mock_balance -mock_names=Repository=BalanceRepository . Repository

type Repository interface {
	// FindUserByID Errors:
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	FindUserByID(userID int64) (*models.Balance, error)
	// CreateTransfer Errors:
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	CreateTransfer(from int64, to int64, amount float64) error
	// CreateAccount Errors:
	// 		app.GeneralError with Errors
	// 			DefaultErrDB
	CreateAccount(userID int64) error
	// AddBalance Errors:
	// 		app.GeneralError with Errors
	// 			DefaultErrDB
	AddBalance(userID int64, amount float64) (float64, error)
}
