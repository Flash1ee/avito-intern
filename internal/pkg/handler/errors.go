package handler

import "github.com/pkg/errors"

var (
	InvalidBody          = errors.New("invalid body")
	InvalidParameters    = errors.New("invalid parameters in query")
	UserNotFound         = errors.New("user not found")
	BDError              = errors.New("can not do bd operation")
	NotEnoughMoney       = errors.New("the user does not have enough funds to perform the operation")
	NotSupportedCurrency = errors.New("not supported currency for convert, check available in API")
	CurrencyConvertError = errors.New("currency convert error, try ordinary request or try later")
)
