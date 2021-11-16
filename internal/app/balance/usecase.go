package balance

import "avito-intern/internal/app/balance/models"

//go:generate mockgen -destination=mocks/usecase.go -package=mock_balance -mock_names=Usecase=Usecase . Usecase

type Usecase interface {
	// UpdateBalance Errors:
	//		balance_usecase.NotEnoughMoney
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	UpdateBalance(userID int64, amount float64, updateType int) (float64, error)

	// GetBalance Errors:
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	GetBalance(userID int64) (float64, error)
	// TransferMoney Errors:
	//		balance_usecase.NotEnoughMoney
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	TransferMoney(senderID int64, receiverID int64, amount float64) (*models.TransferMoney, error)
}
