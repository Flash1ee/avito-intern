package transaction_usecase

import (
	test_data "avito-intern/internal/app/balance/testing"
	"avito-intern/internal/app/transaction"
	"avito-intern/internal/app/transaction/models"
	"avito-intern/internal/app/transaction/transaction_repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteTransactionUsecase struct {
	SuiteUsecase
	uc transaction.Usecase
}

func (s *SuiteTransactionUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewTransactionUsecase(s.MockTransactionRepository)
}

func (s *SuiteTransactionUsecase) TestTransactionUsecase_GetTransactions_OK() {
	b := test_data.TestBalance(s.T())
	paginator := test_data.TestPaginator(s.T())

	var mockRes []models.Transaction

	s.MockTransactionRepository.EXPECT().
		GetTransactions(b.ID, paginator).
		Return(mockRes, nil)

	res, err := s.uc.GetTransactions(b.ID, paginator)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), mockRes, res)
}
func (s *SuiteTransactionUsecase) TestTransactionUsecase_GetTransactions_Error() {
	b := test_data.TestBalance(s.T())
	paginator := test_data.TestPaginator(s.T())

	repoErr := transaction_repository.NotFound

	s.MockTransactionRepository.EXPECT().
		GetTransactions(b.ID, paginator).
		Return(nil, repoErr)

	res, err := s.uc.GetTransactions(b.ID, paginator)
	assert.Equal(s.T(), repoErr, err)
	assert.Nil(s.T(), res)
}
func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteTransactionUsecase))
}
