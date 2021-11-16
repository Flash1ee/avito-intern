package balance_usecase

import (
	mock_balance "avito-intern/internal/app/balance/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"io"
)

type SuiteUsecase struct {
	suite.Suite
	Mock                  *gomock.Controller
	MockBalanceRepository *mock_balance.BalanceRepository

	Logger *logrus.Logger
}

func (s *SuiteUsecase) SetupSuite() {
	s.Mock = gomock.NewController(s.T())
	s.MockBalanceRepository = mock_balance.NewBalanceRepository(s.Mock)
	s.Logger = logrus.New()
	s.Logger.SetOutput(io.Discard)
}

func (s *SuiteUsecase) TearDownSuite() {
	s.Mock.Finish()
}
