package balance

type Repository interface {
	// FindUserByID Errors:
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	FindUserByID(userID int64) (int64, error)
	// GetBalance Errors:
	//		balance_repository.NotFound
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	GetBalance(userID int64) (int64, error)
	// AddTransaction Errors:
	// 		app.GeneralError with Errors
	// 			balance_repository.DefaultErrDB
	AddTransaction(fromID int64, toID int64, amount int64) error
	// CreateAccount Errors:
	// 		app.GeneralError with Errors
	// 			DefaultErrDB
	CreateAccount(userID int64) error
	// AddBalance Errors:
	// 		app.GeneralError with Errors
	// 			DefaultErrDB
	AddBalance(userID int64, amount int64) (int64, error)
}
