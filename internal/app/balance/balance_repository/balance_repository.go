package balance_repository

import (
	"avito-intern/internal/app"
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
func (repo *BalanceRepository) FindUserByID(userID int64) (int64, error) {
	query := "SELECT balance from balance where user_id = $1"
	var balance int64

	if err := repo.conn.QueryRow(query, userID).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app.InvalidInt, NotFound
		}
		return app.InvalidInt, NewDBError(err)
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

// AddTransaction Errors:
// 		app.GeneralError with Errors
// 			DefaultErrDB
func (repo *BalanceRepository) AddTransaction(fromID int64, toID int64, amount int64) error {
	query := "INSERT INTO transactions(from_id, to_id, amount)" +
		"VALUES($1, $2, $3)"

	res := repo.conn.QueryRow(query, fromID, toID, amount)
	if res.Err() != nil {
		return NewDBError(res.Err())
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
