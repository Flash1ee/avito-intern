package balance_repository

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance/models"
	"database/sql"
	"github.com/pkg/errors"
)

type BalanceRepository struct {
	conn *sql.DB
}

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
	query := "SELECT user_id, amount from balance where user_id = $1"
	balance := &models.Balance{}

	if err := repo.conn.QueryRow(query, userID).Scan(&balance.ID, &balance.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		}
		return nil, NewDBError(err)
	}
	return balance, nil
}

// GetBalance Errors:
//		NotFound
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) GetBalance(userID int64) (int64, error) {
	query := "SELECT amount from balance where user_id = $1"

	var resBalance int64
	err := repo.conn.QueryRow(query, userID).Scan(&resBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, NotFound
		}
		return app.InvalidInt, NewDBError(err)
	}
	return resBalance, nil
}

// CreateTransfer Errors:
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) CreateTransfer(from int64, to int64, amount int64) error {
	queryWriteOff := "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"
	queryEnroll := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2"
	queryAddTransaction := "INSERT INTO transactions(from_id, to_id, amount) VALUES($1, $2, $3)"

	transact, err := repo.conn.Begin()
	if err != nil {
		return NewDBError(err)
	}
	_, err = transact.Exec(queryWriteOff, amount, from)
	if err != nil {
		_ = transact.Rollback()
		return NewDBError(err)
	}
	_, err = transact.Exec(queryEnroll, amount, to)
	if err != nil {
		_ = transact.Rollback()
		return NewDBError(err)
	}
	_, err = transact.Exec(queryAddTransaction, from, to, amount)
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
	query := "INSERT INTO balance(user_id) VALUES($1)"

	res := repo.conn.QueryRow(query, userID)
	if res.Err() != nil {
		return NewDBError(res.Err())
	}
	return nil
}

// AddBalance Errors:
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) AddBalance(userID int64, amount int64) (int64, error) {
	query := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2 RETURNING amount"

	var balance int64
	err := repo.conn.QueryRow(query, amount, userID).Scan(&balance)
	if err != nil {
		return app.InvalidInt, NewDBError(err)
	}
	return balance, nil
}
