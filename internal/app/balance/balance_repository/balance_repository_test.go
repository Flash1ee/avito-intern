package balance_repository

import (
	"avito-intern/internal/app"
	test_data "avito-intern/internal/app/balance/testing"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
)

type SuiteBalanceRepository struct {
	Suite
	repo *BalanceRepository
}

func (s *SuiteBalanceRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewBalanceRepository(s.DB)
}

func (s *SuiteBalanceRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}
func (s *SuiteBalanceRepository) TestBalanceRepository_FindUserByID_ERROR_DB() {
	query := "SELECT user_id, amount from balance where user_id = $1"
	userID := int64(1)
	expError := DefaultErrDB
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID).WillReturnError(DefaultErrDB)

	user, err := s.repo.FindUserByID(userID)
	assert.Nil(s.T(), user)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestBalanceRepository_FindUserByID_NotFound() {
	query := "SELECT user_id, amount from balance where user_id = $1"
	userID := int64(1)
	expError := NotFound

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID).WillReturnError(sql.ErrNoRows)

	user, err := s.repo.FindUserByID(userID)
	assert.Nil(s.T(), user)
	assert.Equal(s.T(), expError, err)
}

func (s *SuiteBalanceRepository) TestBalanceRepository_FindUserByID_OK() {
	query := "SELECT user_id, amount from balance where user_id = $1"
	userID := int64(1)
	expBalance := test_data.TestBalance(s.T())

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"user_id", "amount"}).
			AddRow(expBalance.ID, expBalance.Amount))

	balance, err := s.repo.FindUserByID(userID)

	assert.Equal(s.T(), expBalance, balance)
	assert.Nil(s.T(), err)
}
func (s *SuiteBalanceRepository) TestCreateAccount_ERROR_DB() {
	query := "INSERT INTO balance(user_id) VALUES($1)"

	userID := int64(1)
	expError := DefaultErrDB
	s.Mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(userID).WillReturnError(DefaultErrDB)

	err := s.repo.CreateAccount(userID)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateAccount_OK() {
	query := "INSERT INTO balance(user_id) VALUES($1)"

	userID := int64(1)
	s.Mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.CreateAccount(userID)
	assert.NoError(s.T(), err)
}
func (s *SuiteBalanceRepository) TestAddBalance_ERROR_DB() {
	query := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2 RETURNING amount"
	expErr := DefaultErrDB

	testB := test_data.TestBalance(s.T())
	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(testB.Amount, testB.ID).
		WillReturnError(expErr)

	res, err := s.repo.AddBalance(testB.ID, testB.Amount)
	assert.Equal(s.T(), NewDBError(expErr), err)
	assert.Equal(s.T(), app.InvalidFloat, res)
}
func (s *SuiteBalanceRepository) TestAddBalance_OK() {
	query := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2 RETURNING amount"
	testB := test_data.TestBalance(s.T())

	expRes := testB.Amount * 2

	s.Mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(testB.Amount, testB.ID).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).
			AddRow(expRes))

	res, err := s.repo.AddBalance(testB.ID, testB.Amount)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_TransactionError() {
	expError := DefaultErrDB
	testT := test_data.TestTransaction(s.T())
	s.Mock.ExpectBegin().WillReturnError(DefaultErrDB)

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, 10)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_WriteOffError() {
	queryWriteOff := "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"

	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOff)).WithArgs(amount, testT.SenderID).
		WillReturnError(expError)
	s.Mock.ExpectRollback()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_EnrollError() {
	queryWriteOff := "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"
	queryEnroll := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2"

	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOff)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnroll)).WithArgs(amount, testT.ReceiverID).
		WillReturnError(expError)
	s.Mock.ExpectRollback()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_AddTransactionError() {
	queryWriteOff := "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"
	queryEnroll := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2"
	queryAddTransaction := "INSERT INTO transactions(from_id, to_id, amount) VALUES($1, $2, $3)"

	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOff)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnroll)).WithArgs(amount, testT.ReceiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(testT.SenderID, testT.ReceiverID, amount).
		WillReturnError(expError)
	s.Mock.ExpectRollback()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_CommitError() {
	queryWriteOff := "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"
	queryEnroll := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2"
	queryAddTransaction := "INSERT INTO transactions(from_id, to_id, amount) VALUES($1, $2, $3)"

	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOff)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnroll)).WithArgs(amount, testT.ReceiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(testT.SenderID, testT.ReceiverID, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit().WillReturnError(expError)

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_OK() {
	queryWriteOff := "UPDATE balance SET amount = amount - $1 WHERE user_id = $2"
	queryEnroll := "UPDATE balance SET amount = amount + $1 WHERE user_id = $2"
	queryAddTransaction := "INSERT INTO transactions(from_id, to_id, amount) VALUES($1, $2, $3)"

	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOff)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnroll)).WithArgs(amount, testT.ReceiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(testT.SenderID, testT.ReceiverID, amount).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Nil(s.T(), err)
}
func TestBalanceRepository(t *testing.T) {
	suite.Run(t, new(SuiteBalanceRepository))
}
