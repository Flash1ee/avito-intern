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
		if err != balance_repository.NotFound || updateType != models.ADD_BALANCE {
			return app.InvalidInt, err
		}
	}
	if balance == nil {
		err = uc.repo.CreateAccount(userID)
		if err != nil {
			return app.InvalidInt, err
		}
	}
	if updateType == models.DIFF_BALANCE {
		if balance.Amount < amount {
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

// TransferMoney Errors:
//		NotEnoughMoney
//		balance_repository.NotFound
// 		app.GeneralError with Errors
// 			balance_repository.DefaultErrDB
func (uc *BalanceUsecase) TransferMoney(senderID int64, receiverID int64, amount int64) (*models.TransferMoney, error) {
	sender, err := uc.repo.FindUserByID(senderID)
	if err != nil {
		return nil, err
	}
	if sender.Amount < amount {
		return nil, NotEnoughMoney
	}
	receiver, err := uc.repo.FindUserByID(receiverID)
	if err != nil {
		return nil, err
	}
	err = uc.repo.CreateTransfer(senderID, receiverID, amount)
	if err != nil {
		return nil, err
	}
	return &models.TransferMoney{
		SenderID:        senderID,
		ReceiverID:      receiverID,
		SenderBalance:   sender.Amount - amount,
		ReceiverBalance: receiver.Amount + amount,
	}, nil
}
