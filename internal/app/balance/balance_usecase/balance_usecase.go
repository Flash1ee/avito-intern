package balance_usecase

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance"
	"avito-intern/internal/app/balance/balance_repository"
	"avito-intern/internal/app/balance/models"
)

type BalanceUsecase struct {
	repo balance.Repository
}

func NewBalanceUsecase(repo balance.Repository) *BalanceUsecase {
	return &BalanceUsecase{
		repo: repo,
	}
}

// GetBalance Errors:
//		balance_repository.NotFound
// 		app.GeneralError with Errors
// 			balance_repository.DefaultErrDB
func (uc *BalanceUsecase) GetBalance(userID int64) (float64, error) {
	user, err := uc.repo.FindUserByID(userID)
	if err != nil {
		return app.InvalidFloat, err
	}
	return user.Amount, nil
}

// UpdateBalance Errors:
//		NotEnoughMoney
//		balance_repository.NotFound
// 		app.GeneralError with Errors
// 			balance_repository.DefaultErrDB
func (uc *BalanceUsecase) UpdateBalance(userID int64, amount float64, updateType int) (float64, error) {
	b, err := uc.repo.FindUserByID(userID)
	if err != nil {
		if err != balance_repository.NotFound || updateType != models.ADD_BALANCE {
			return app.InvalidFloat, err
		}
	}
	if b == nil {
		err = uc.repo.CreateAccount(userID)
		if err != nil {
			return app.InvalidFloat, err
		}
	}
	if updateType == models.DIFF_BALANCE {
		if b.Amount < amount {
			return app.InvalidFloat, NotEnoughMoney
		}
		amount *= -1
	}
	newBalance, err := uc.repo.AddBalance(userID, amount)

	if err != nil {
		return app.InvalidFloat, err
	}
	return newBalance, nil
}

// TransferMoney Errors:
//		NotEnoughMoney
//		balance_repository.NotFound
// 		app.GeneralError with Errors
// 			balance_repository.DefaultErrDB
func (uc *BalanceUsecase) TransferMoney(senderID int64, receiverID int64, amount float64) (*models.TransferMoney, error) {
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
