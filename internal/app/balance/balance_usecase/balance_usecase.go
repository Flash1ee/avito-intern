package balance_usecase

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance/balance_repository"
	"avito-intern/internal/app/balance/models"
)

type BalanceUsecase struct {
	repo *balance_repository.BalanceRepository
}

func NewBalanceUsecase(repo *balance_repository.BalanceRepository) *BalanceUsecase {
	return &BalanceUsecase{
		repo: repo,
	}
}

// GetBalance Errors:
//		balance_repository.NotFound
// 		app.GeneralError with Errors
// 			balance_repository.DefaultErrDB
func (uc *BalanceUsecase) GetBalance(userID int64) (int64, error) {
	balance, err := uc.repo.GetBalance(userID)
	if err != nil {
		return app.InvalidInt, err
	}
	return balance, nil
}

// UpdateBalance Errors:
//		NotEnoughMoney
//		balance_repository.NotFound
// 		app.GeneralError with Errors
// 			balance_repository.DefaultErrDB
func (uc *BalanceUsecase) UpdateBalance(userID int64, amount int64, updateType int) (int64, error) {
	balance, err := uc.repo.FindUserByID(userID)
	if err != nil {
		return app.InvalidInt, err
	}
	if balance == app.InvalidInt {
		err = uc.repo.CreateAccount(userID)
		if err != nil {
			return app.InvalidInt, err
		}
	}
	if updateType == models.DIFF_BALANCE {
		if balance < amount {
			return app.InvalidInt, NotEnoughMoney
		}
		amount *= -1
	}
	newBalance, err := uc.repo.AddBalance(userID, amount)

	if err != nil {
		return app.InvalidInt, err
	}
	return newBalance, nil
}
