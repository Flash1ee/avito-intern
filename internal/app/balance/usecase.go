package balance

type Usecase interface {
	// UpdateBalance Errors:
	//		balance_usecase.NotEnoughMoney
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	UpdateBalance(userID int64, amount int64, updateType int) (int64, error)

	// GetBalance Errors:
	//		NotFound
	// 		app.GeneralError with Errors
	// 			DefaultErrDB
	GetBalance(userID int64) (int64, error)
}
