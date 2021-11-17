package transaction_handler

import (
	mock_transaction "avito-intern/internal/app/transaction/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
)

type TestTable struct {
	Name              string
	Data              interface{}
	ExpectedMockTimes int
	ExpectedCode      int
}

type SuiteHandler struct {
	suite.Suite
	Mock                   *gomock.Controller
	MockTransactionUsecase *mock_transaction.TransactionUsecase
	Logger                 *logrus.Logger
}

func (s *SuiteHandler) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockTransactionUsecase = mock_transaction.NewTransactionUsecase(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteHandler) TearDownSuite() {
	s.Mock.Finish()
}
