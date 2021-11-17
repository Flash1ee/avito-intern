package transaction_usecase

import (
	mock_transaction "avito-intern/internal/app/transaction/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
)

type SuiteUsecase struct {
	suite.Suite
	Mock                      *gomock.Controller
	MockTransactionRepository *mock_transaction.TransactionRepository
	Logger                    *logrus.Logger
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockTransactionRepository = mock_transaction.NewTransactionRepository(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
