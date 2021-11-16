package balance_handler

import (
	"avito-intern/internal/app/balance/balance_repository"
	"avito-intern/internal/app/balance/balance_usecase"
	"avito-intern/internal/pkg/handler"
	"github.com/sirupsen/logrus"
	"net/http"
)

var CodeByErrorGetBalance = handler.CodeMap{
	balance_repository.NotFound:     {http.StatusNotFound, handler.UserNotFound, logrus.WarnLevel},
	balance_repository.DefaultErrDB: {http.StatusInternalServerError, handler.BDError, logrus.ErrorLevel},
}
var CodeByErrorUpdateBalanceHandler = handler.CodeMap{
	balance_repository.NotFound:     {http.StatusNotFound, handler.UserNotFound, logrus.WarnLevel},
	balance_repository.DefaultErrDB: {http.StatusInternalServerError, handler.BDError, logrus.ErrorLevel},
	balance_usecase.NotEnoughMoney:  {http.StatusUnprocessableEntity, handler.NotEnoughMoney, logrus.ErrorLevel},
}
var CodeByErrorTransferHandler = handler.CodeMap{
	balance_repository.NotFound:     {http.StatusNotFound, handler.UserNotFound, logrus.WarnLevel},
	balance_repository.DefaultErrDB: {http.StatusInternalServerError, handler.BDError, logrus.ErrorLevel},
	balance_usecase.NotEnoughMoney:  {http.StatusUnprocessableEntity, handler.NotEnoughMoney, logrus.ErrorLevel},
}
