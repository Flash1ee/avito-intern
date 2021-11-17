package balance_repository

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance/models"
	transaction_models "avito-intern/internal/app/models"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

type BalanceRepository struct {
	conn *sql.DB
}

const (
	queryFindUserById    = "SELECT user_id, amount from balance where user_id = $1"
	queryWriteOffBalance = "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"
	queryEnrollBalance   = "UPDATE balance SET amount = amount + $1 WHERE user_id = $2"
	queryAddTransaction  = "INSERT INTO transactions(type, sender_id, receiver_id, amount, description) VALUES($1, $2, $3, $4, $5)"
	queryCreateAccount   = "INSERT INTO balance(user_id) VALUES($1)"
	queryAddBalance      = "UPDATE balance SET amount = amount + $1 WHERE user_id = $2 RETURNING amount"
)
const (
	addBalanceDescriptions = "user <id = %d> %s account"
	transferDescription    = "user <id = %d> send money to user <id = %d>"
)

func NewBalanceRepository(conn *sql.DB) *BalanceRepository {
	return &BalanceRepository{
		conn: conn,
	}
}

// FindUserByID Errors:
//		NotFound
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) FindUserByID(userID int64) (*models.Balance, error) {
	balance := &models.Balance{}

	if err := repo.conn.QueryRow(queryFindUserById, userID).Scan(&balance.ID, &balance.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		}
		return nil, NewDBError(err)
	}
	return balance, nil
}

// CreateTransfer Errors:
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) CreateTransfer(from int64, to int64, amount float64) error {
	operationType := transaction_models.TextTransactionToType[transaction_models.TRANSFER]
	description := fmt.Sprintf(transferDescription, from, to)

	transact, err := repo.conn.Begin()
	if err != nil {
		return NewDBError(err)
	}
	_, err = transact.Exec(queryWriteOffBalance, amount, from)
	if err != nil {
		_ = transact.Rollback()
		return NewDBError(err)
	}
	_, err = transact.Exec(queryEnrollBalance, amount, to)
	if err != nil {
		_ = transact.Rollback()
		return NewDBError(err)
	}
	_, err = transact.Exec(queryAddTransaction, operationType, from, to, amount, description)
	if err != nil {
		_ = transact.Rollback()
		return NewDBError(err)
	}

	if err = transact.Commit(); err != nil {
		return NewDBError(err)
	}

	return nil
}

// CreateAccount Errors:
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) CreateAccount(userID int64) error {
	_, err := repo.conn.Exec(queryCreateAccount, userID)
	if err != nil {
		return NewDBError(err)
	}

	return nil
}

// AddBalance Errors:
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) AddBalance(userID int64, amount float64) (float64, error) {
	var description string
	operation := transaction_models.TextTransactionToType[transaction_models.REFILL]
	if amount < 0 {
		description = fmt.Sprintf(
			addBalanceDescriptions, userID, transaction_models.TextTransactionToType[transaction_models.WRITE_OFF])
		operation = transaction_models.TextTransactionToType[transaction_models.WRITE_OFF]
	} else {
		description = fmt.Sprintf(
			addBalanceDescriptions, userID, transaction_models.TextTransactionToType[transaction_models.REFILL])
	}

	var balance float64
	err := repo.conn.QueryRow(queryAddBalance, amount, userID).Scan(&balance)
	if err != nil {
		return app.InvalidFloat, NewDBError(err)
	}

	_, err = repo.conn.Exec(queryAddTransaction, operation, userID, userID, amount, description)
	if err != nil {
		return app.InvalidFloat, NewDBError(err)
	}

	return balance, nil
}
