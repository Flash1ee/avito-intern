package balance_repository

import (
	"avito-intern/internal/app"
	test_data "avito-intern/internal/app/balance/testing"
	transaction_models "avito-intern/internal/app/models"
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
	userID := int64(1)
	expError := DefaultErrDB
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryFindUserById)).
		WithArgs(userID).WillReturnError(DefaultErrDB)

	user, err := s.repo.FindUserByID(userID)
	assert.Nil(s.T(), user)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestBalanceRepository_FindUserByID_NotFound() {
	userID := int64(1)
	expError := NotFound

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryFindUserById)).
		WithArgs(userID).WillReturnError(sql.ErrNoRows)

	user, err := s.repo.FindUserByID(userID)
	assert.Nil(s.T(), user)
	assert.Equal(s.T(), expError, err)
}

func (s *SuiteBalanceRepository) TestBalanceRepository_FindUserByID_OK() {
	userID := int64(1)
	expBalance := test_data.TestBalance(s.T())

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryFindUserById)).
		WithArgs(userID).WillReturnRows(
		sqlmock.NewRows([]string{"user_id", "amount"}).
			AddRow(expBalance.ID, expBalance.Amount))

	balance, err := s.repo.FindUserByID(userID)

	assert.Equal(s.T(), expBalance, balance)
	assert.Nil(s.T(), err)
}
func (s *SuiteBalanceRepository) TestCreateAccount_ERROR_DB() {
	userID := int64(1)
	expError := DefaultErrDB
	s.Mock.ExpectExec(regexp.QuoteMeta(queryCreateAccount)).
		WithArgs(userID).WillReturnError(DefaultErrDB)

	err := s.repo.CreateAccount(userID)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateAccount_OK() {
	userID := int64(1)
	s.Mock.ExpectExec(regexp.QuoteMeta(queryCreateAccount)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.CreateAccount(userID)
	assert.NoError(s.T(), err)
}
func (s *SuiteBalanceRepository) TestAddBalance_ERROR_DB() {
	expErr := DefaultErrDB

	testB := test_data.TestBalance(s.T())
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddBalance)).
		WithArgs(testB.Amount, testB.ID).
		WillReturnError(expErr)

	res, err := s.repo.AddBalance(testB.ID, testB.Amount)
	assert.Equal(s.T(), NewDBError(expErr), err)
	assert.Equal(s.T(), app.InvalidFloat, res)
}
func (s *SuiteBalanceRepository) TestAddBalance_WriteTransactionError() {
	testB := test_data.TestBalance(s.T())
	operationType := transaction_models.TextTransactionToType[transaction_models.REFILL]

	testDecr := test_data.TestAddBalanceDescription(testB.ID, operationType, s.T())

	expRes := testB.Amount * 2
	sqlErr := sql.ErrTxDone
	expErr := NewDBError(sqlErr)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddBalance)).
		WithArgs(testB.Amount, testB.ID).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).
			AddRow(expRes))

	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(operationType, testB.ID, testB.ID, testB.Amount, testDecr).
		WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(sqlErr)

	res, err := s.repo.AddBalance(testB.ID, testB.Amount)
	assert.Equal(s.T(), expErr, err)
	assert.Equal(s.T(), app.InvalidFloat, res)
}
func (s *SuiteBalanceRepository) TestAddBalance_OK() {
	testB := test_data.TestBalance(s.T())
	operationType := transaction_models.TextTransactionToType[transaction_models.REFILL]

	testDecr := test_data.TestAddBalanceDescription(testB.ID, operationType, s.T())

	expRes := testB.Amount * 2

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddBalance)).
		WithArgs(testB.Amount, testB.ID).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).
			AddRow(expRes))

	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(operationType, testB.ID, testB.ID, testB.Amount, testDecr).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := s.repo.AddBalance(testB.ID, testB.Amount)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceRepository) TestAddBalance_OK_WriteOff() {
	testB := test_data.TestBalance(s.T())
	operationType := transaction_models.TextTransactionToType[transaction_models.WRITE_OFF]

	testDecr := test_data.TestAddBalanceDescription(testB.ID, operationType, s.T())

	expRes := testB.Amount / 2
	writeOffMoney := (testB.Amount / 2) * -1

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryAddBalance)).
		WithArgs(writeOffMoney, testB.ID).
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).
			AddRow(expRes))

	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(operationType, testB.ID, testB.ID, writeOffMoney, testDecr).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := s.repo.AddBalance(testB.ID, writeOffMoney)
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
	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOffBalance)).WithArgs(amount, testT.SenderID).
		WillReturnError(expError)
	s.Mock.ExpectRollback()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_EnrollError() {
	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOffBalance)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnrollBalance)).WithArgs(amount, testT.ReceiverID).
		WillReturnError(expError)
	s.Mock.ExpectRollback()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_AddTransactionError() {
	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	transferDecr := test_data.TestTransferDescription(testT.SenderID, testT.ReceiverID, s.T())
	operType := transaction_models.TextTransactionToType[transaction_models.TRANSFER]
	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOffBalance)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnrollBalance)).WithArgs(amount, testT.ReceiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(operType, testT.SenderID, testT.ReceiverID, amount, transferDecr).
		WillReturnError(expError)

	s.Mock.ExpectRollback()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_CommitError() {
	expError := DefaultErrDB
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	transferDecr := test_data.TestTransferDescription(testT.SenderID, testT.ReceiverID, s.T())
	operType := transaction_models.TextTransactionToType[transaction_models.TRANSFER]

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOffBalance)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnrollBalance)).WithArgs(amount, testT.ReceiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(operType, testT.SenderID, testT.ReceiverID, amount, transferDecr).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit().WillReturnError(expError)

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Equal(s.T(), NewDBError(expError), err)
}
func (s *SuiteBalanceRepository) TestCreateTransfer_OK() {
	amount := float64(10)
	testT := test_data.TestTransaction(s.T())

	transferDecr := test_data.TestTransferDescription(testT.SenderID, testT.ReceiverID, s.T())
	operType := transaction_models.TextTransactionToType[transaction_models.TRANSFER]

	s.Mock.ExpectBegin()
	s.Mock.ExpectExec(regexp.QuoteMeta(queryWriteOffBalance)).WithArgs(amount, testT.SenderID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryEnrollBalance)).WithArgs(amount, testT.ReceiverID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectExec(regexp.QuoteMeta(queryAddTransaction)).
		WithArgs(operType, testT.SenderID, testT.ReceiverID, amount, transferDecr).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.Mock.ExpectCommit()

	err := s.repo.CreateTransfer(testT.SenderID, testT.ReceiverID, amount)
	assert.Nil(s.T(), err)
}
func TestBalanceRepository(t *testing.T) {
	suite.Run(t, new(SuiteBalanceRepository))
}
