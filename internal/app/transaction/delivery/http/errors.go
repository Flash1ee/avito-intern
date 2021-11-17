package transaction_handler

import "github.com/pkg/errors"

var (
	InvalidQueryPageParam      = errors.New("invalid query param page")
	InvalidQueryCountParam     = errors.New("invalid query param count")
	InvalidQuerySortParam      = errors.New("invalid query param sort")
	InvalidQueryDirectionParam = errors.New("invalid query param direction")
	DirectionMustUsageWithSort = errors.New("for the sorting direction, you need a " +
		"sorting condition - the sort field")
	TransactionsNotFound = errors.New("transactions not found")
)
