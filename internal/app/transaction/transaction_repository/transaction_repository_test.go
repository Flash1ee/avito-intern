package transaction_repository

import (
	"avito-intern/internal/app"
	test_data "avito-intern/internal/app/balance/testing"
	transaction_models "avito-intern/internal/app/models"
	"avito-intern/internal/app/transaction/models"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"regexp"
	"testing"
	"time"
)

type SuiteTransactionRepository struct {
	Suite
	repo *TransactionRepository
}

func (s *SuiteTransactionRepository) SetupSuite() {
	s.InitBD()
	s.repo = NewTransactionRepository(s.DB)
}

func (s *SuiteTransactionRepository) AfterTest(_, _ string) {
	require.NoError(s.T(), s.Mock.ExpectationsWereMet())
}

func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_NotFound() {
	testPaginator := test_data.TestPaginator(s.T())
	userID := int64(1)
	expError := NotFound

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(0))
	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expError, err)
}
func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_CountSelectError() {
	testPaginator := test_data.TestPaginator(s.T())
	userID := int64(1)
	sqlError := sqlmock.ErrCancelled
	expError := NewDBError(sqlError)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnError(sqlError)
	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expError, err)
}
func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_SelectError() {
	testPaginator := test_data.TestPaginator(s.T())
	userID := int64(1)
	sqlErr := pq.ErrNotSupported
	expErr := NewDBError(sqlErr)

	expQuerySelect := querySelectTransactions +
		fmt.Sprintf(defineQueryPagination, testPaginator.Count, testPaginator.Page-1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(1))
	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuerySelect)).
		WithArgs(userID).WillReturnError(sqlErr)

	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_ScanError() {
	testPaginator := test_data.TestPaginator(s.T())
	userID := int64(1)
	sqlErr := pq.ErrNotSupported
	expErr := NewDBError(sqlErr)

	expQuerySelect := querySelectTransactions +
		fmt.Sprintf(defineQueryPagination, testPaginator.Count, testPaginator.Page-1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(2))
	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuerySelect)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"type", "sender_id", "receiver_id", "amount", "created_at", "description"}).
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, time.Now(), "description").
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], "invalid id", userID, 100, time.Now(), "description")).
		RowsWillBeClosed()
	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr.Err, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_RowsErr() {
	testPaginator := test_data.TestPaginator(s.T())
	userID := int64(1)
	sqlErr := pq.ErrNotSupported
	expErr := NewDBError(sqlErr)

	expQuerySelect := querySelectTransactions +
		fmt.Sprintf(defineQueryPagination, testPaginator.Count, testPaginator.Page-1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(2))
	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuerySelect)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"type", "sender_id", "receiver_id", "amount", "created_at", "description"}).
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, time.Now(), "description").
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, time.Now(), "description")).
		WillReturnError(sqlErr)

	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr.Err, errors.Cause(err).(*app.GeneralError).Err)
}
func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_OK() {
	testPaginator := test_data.TestPaginator(s.T())
	userID := int64(1)

	expQuerySelect := querySelectTransactions +
		fmt.Sprintf(defineQueryPagination, testPaginator.Count, testPaginator.Page-1)

	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(2))
	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuerySelect)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"type", "sender_id", "receiver_id", "amount", "created_at", "description"}).
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, time.Now(), "description").
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, time.Now(), "description")).
		RowsWillBeClosed()

	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
}
func (s *SuiteTransactionRepository) TestTransactionRepository_GetTransactions_OK_WithOrder() {
	testPaginator := test_data.TestPaginator(s.T())
	testPaginator.SortField = models.DATE
	testPaginator.SortDirection = models.ASC
	userID := int64(1)

	expQuerySelect := querySelectTransactions + fmt.Sprintf(defineQueryOrder, models.TransactionQueryParams[testPaginator.SortField],
		models.TransactionQueryParams[testPaginator.SortDirection])

	expQuerySelect +=
		fmt.Sprintf(defineQueryPagination, testPaginator.Count, testPaginator.Page-1)
	dateOne := time.Now()
	dateTwo := dateOne.Add(time.Second)
	s.Mock.ExpectQuery(regexp.QuoteMeta(queryCountTransactions)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"cnt"}).AddRow(2))
	s.Mock.ExpectQuery(regexp.QuoteMeta(expQuerySelect)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"type", "sender_id", "receiver_id", "amount", "created_at", "description"}).
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, dateOne, "description").
			AddRow(transaction_models.TextTransactionToType[transaction_models.WRITE_OFF], userID, userID, 100, dateTwo, "description")).
		RowsWillBeClosed()

	res, err := s.repo.GetTransactions(userID, testPaginator)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), res)
	assert.Equal(s.T(), dateOne, res[0].CreatedAt)
	assert.Equal(s.T(), dateTwo, res[1].CreatedAt)

}

func TestBalanceRepository(t *testing.T) {
	suite.Run(t, new(SuiteTransactionRepository))
}
