package transaction_repository

import (
	"avito-intern/internal/app/transaction/models"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

type TransactionRepository struct {
	conn *sql.DB
}

func NewTransactionRepository(conn *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		conn: conn,
	}
}

// GetTransactions Errors:
//		NotFound
//		app.GeneralError with Errors:
//			DefaultErrDB
func (repo *TransactionRepository) GetTransactions(userID int64, paginator *models.Paginator) (
	[]models.Transaction, error) {

	queryCount := "SELECT count(*) FROM transactions where sender_id = $1"
	querySelect := "SELECT type, sender_id, receiver_id, amount, created_at, description FROM transactions where sender_id = $1 "

	defineOrder := "ORDER BY %s %s "
	queryPagination := "LIMIT %d OFFSET %d"

	count := 0
	if err := repo.conn.QueryRow(queryCount, userID).Scan(&count); err != nil {
		return nil, NewDBError(err)
	}
	if count == 0 {
		return nil, NotFound
	}

	if paginator.SortField != models.NO_ORDER {
		querySelect += fmt.Sprintf(defineOrder, models.TransactionQueryParams[paginator.SortField],
			models.TransactionQueryParams[paginator.SortDirection])
	}
	querySelect += fmt.Sprintf(queryPagination, paginator.Count, paginator.Page-1)

	rows, err := repo.conn.Query(querySelect, userID)
	if err != nil {
		return nil, NewDBError(err)
	}

	transactionRes := make([]models.Transaction, 0, count)

	for rows.Next() {
		cur := models.Transaction{}
		if err = rows.Scan(&cur.Type, &cur.UserID, &cur.ReceiverID, &cur.Amount,
			&cur.CreatedAt, &cur.Description); err != nil {
			_ = rows.Close()
			return nil, NewDBError(errors.Wrapf(err, "GetTransactions"+
				"invalid data in db: table transactions"))
		}
		transactionRes = append(transactionRes, cur)
	}

	if err = rows.Err(); err != nil {
		return nil, NewDBError(err)
	}

	return transactionRes, nil
}
