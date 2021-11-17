package transaction_handler

import (
	"avito-intern/internal/app/transaction/transaction_repository"
	"avito-intern/internal/pkg/handler"
	"github.com/sirupsen/logrus"
	"net/http"
)

// pagination constants
const (
	PAGE      = "page"
	COUNT     = "count"
	DIRECTION = "direction"
	SORT      = "sort"
	DATE      = "date"
	SUM       = "sum"
	ASC       = "asc"
	DESC      = "desc"
)

var codeByErrorGetTransactions = handler.CodeMap{
	transaction_repository.DefaultErrDB: {http.StatusInternalServerError, handler.BDError, logrus.ErrorLevel},
}
