package transaction_usecase

import (
	mock_transaction "avito-intern/internal/app/transaction/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type SuiteUsecase struct {
	suite.Suite
	Mock                      *gomock.Controller
	MockTransactionRepository *mock_transaction.TransactionRepository
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockTransactionRepository = mock_transaction.NewTransactionRepository(s.Mock)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
