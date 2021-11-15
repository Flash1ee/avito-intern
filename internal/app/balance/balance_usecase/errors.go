package balance_usecase

import "github.com/pkg/errors"

var (
	NotEnoughMoney = errors.New("the user does not have enough funds to perform the operation")
)
