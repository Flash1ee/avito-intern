package balance_usecase

import (
	"avito-intern/internal/app/balance"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SuiteBalanceUsecase struct {
	SuiteUsecase
	uc balance.Usecase
}

func (s *SuiteBalanceUsecase) SetupSuite() {
	s.SuiteUsecase.SetupSuite()
	s.uc = NewBalanceUsecase(s.MockBalanceRepository)
}

func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteBalanceUsecase))
}
