package handler

import "github.com/pkg/errors"

var (
	UserNotFound   = errors.New("user not found")
	BDError        = errors.New("can not do bd operation")
	NotEnoughMoney = errors.New("the user does not have enough funds to perform the operation")
)
