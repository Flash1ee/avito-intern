package balance_usecase

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance"
	"avito-intern/internal/app/balance/balance_repository"
	test_data "avito-intern/internal/app/balance/balance_repository/testing"
	"avito-intern/internal/app/balance/models"
	"database/sql"
	"github.com/stretchr/testify/assert"
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

func (s *SuiteBalanceUsecase) TestBalanceUsecase_GetBalance_NotFound() {
	b := test_data.TestBalance(s.T())
	expErr := balance_repository.NotFound

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, balance_repository.NotFound)
	res, err := s.uc.GetBalance(b.ID)
	assert.Equal(s.T(), app.InvalidFloat, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_GetBalance_InteranalError() {
	b := test_data.TestBalance(s.T())
	errRepo := sql.ErrTxDone
	expErr := balance_repository.NewDBError(errRepo)

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, expErr)
	res, err := s.uc.GetBalance(b.ID)
	assert.Equal(s.T(), app.InvalidFloat, res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_GetBalance_OK() {
	b := test_data.TestBalance(s.T())

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(b, nil)

	res, err := s.uc.GetBalance(b.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), b.Amount, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_NotFound_DiffBalance() {
	b := test_data.TestBalance(s.T())

	repoErr := balance_repository.NotFound
	operation := models.DIFF_BALANCE

	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, repoErr)

	res, err := s.uc.UpdateBalance(b.ID, b.Amount, operation)

	assert.Equal(s.T(), repoErr, err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_InternalError_DiffBalance() {
	b := test_data.TestBalance(s.T())

	repoErr := balance_repository.NewDBError(sql.ErrTxDone)
	operation := models.DIFF_BALANCE

	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, repoErr)

	res, err := s.uc.UpdateBalance(b.ID, b.Amount, operation)

	assert.Equal(s.T(), repoErr, err)
	assert.Equal(s.T(), expRes, res)

}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_InternalError_AddBalance() {
	b := test_data.TestBalance(s.T())

	repoErr := balance_repository.NewDBError(sql.ErrTxDone)
	operation := models.ADD_BALANCE

	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, repoErr)

	res, err := s.uc.UpdateBalance(b.ID, b.Amount, operation)

	assert.Equal(s.T(), repoErr, err)
	assert.Equal(s.T(), expRes, res)

}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_NotFound_AddBalance_CreateAccError() {
	b := test_data.TestBalance(s.T())

	repoErr := balance_repository.NotFound
	operation := models.ADD_BALANCE

	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, repoErr)
	s.MockBalanceRepository.EXPECT().
		CreateAccount(b.ID).
		Return(repoErr)

	res, err := s.uc.UpdateBalance(b.ID, b.Amount, operation)

	assert.Equal(s.T(), repoErr, err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_NotFound_AddBalance_AddBalanceError() {
	b := test_data.TestBalance(s.T())

	repoErr := balance_repository.NotFound
	operation := models.ADD_BALANCE

	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, repoErr)
	s.MockBalanceRepository.EXPECT().
		CreateAccount(b.ID).
		Return(nil)
	s.MockBalanceRepository.EXPECT().
		AddBalance(b.ID, b.Amount).
		Return(app.InvalidFloat, repoErr)

	res, err := s.uc.UpdateBalance(b.ID, b.Amount, operation)

	assert.Equal(s.T(), repoErr, err)
	assert.Equal(s.T(), expRes, res)

}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_NotFound_AddBalance_OK() {
	b := test_data.TestBalance(s.T())

	operation := models.ADD_BALANCE
	repoErr := balance_repository.NotFound
	addMoney := 100.2
	expRes := b.Amount + addMoney

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, repoErr)
	s.MockBalanceRepository.EXPECT().
		CreateAccount(b.ID).
		Return(nil)
	s.MockBalanceRepository.EXPECT().
		AddBalance(b.ID, b.Amount).
		Return(expRes, nil)

	res, err := s.uc.UpdateBalance(b.ID, b.Amount, operation)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expRes, res)

}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_UserFound_DiffBalance_NotEnoughMoney() {
	b := test_data.TestBalance(s.T())

	operation := models.DIFF_BALANCE
	diffMoney := b.Amount * 2

	expErr := NotEnoughMoney
	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(b, nil)

	res, err := s.uc.UpdateBalance(b.ID, diffMoney, operation)

	assert.Equal(s.T(), expErr, err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_UserFound_DiffBalance_AddBalanceError() {
	b := test_data.TestBalance(s.T())

	operation := models.DIFF_BALANCE
	diffMoney := b.Amount

	expErr := balance_repository.NewDBError(balance_repository.DefaultErrDB)
	expRes := app.InvalidFloat

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(b, nil)
	s.MockBalanceRepository.EXPECT().
		AddBalance(b.ID, -1*diffMoney).
		Return(expRes, expErr)

	res, err := s.uc.UpdateBalance(b.ID, diffMoney, operation)

	assert.Equal(s.T(), expErr, err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_UserFound_DiffBalance_OK() {
	b := test_data.TestBalance(s.T())

	operation := models.DIFF_BALANCE
	diffMoney := b.Amount / 2

	expRes := b.Amount - diffMoney

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(b, nil)
	s.MockBalanceRepository.EXPECT().
		AddBalance(b.ID, -1*diffMoney).
		Return(expRes, nil)

	res, err := s.uc.UpdateBalance(b.ID, diffMoney, operation)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_UserFound_AddBalance_OK() {
	b := test_data.TestBalance(s.T())

	operation := models.ADD_BALANCE
	diffMoney := b.Amount / 2

	expRes := b.Amount + diffMoney

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(b, nil)
	s.MockBalanceRepository.EXPECT().
		AddBalance(b.ID, diffMoney).
		Return(expRes, nil)

	res, err := s.uc.UpdateBalance(b.ID, diffMoney, operation)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestBalanceUsecase_UpdateBalance_NewAccount_AddBalance_OK() {
	b := test_data.TestBalance(s.T())

	operation := models.ADD_BALANCE
	diffMoney := b.Amount / 2

	expRes := b.Amount + diffMoney

	s.MockBalanceRepository.EXPECT().FindUserByID(b.ID).
		Return(nil, nil)
	s.MockBalanceRepository.EXPECT().
		CreateAccount(b.ID).
		Return(nil)
	s.MockBalanceRepository.EXPECT().
		AddBalance(b.ID, diffMoney).
		Return(expRes, nil)

	res, err := s.uc.UpdateBalance(b.ID, diffMoney, operation)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), expRes, res)
}
func (s *SuiteBalanceUsecase) TestTransferMoney_SenderNotFound() {
	sender := test_data.TestBalance(s.T())
	expErr := balance_repository.NotFound

	s.MockBalanceRepository.EXPECT().
		FindUserByID(sender.ID).
		Return(nil, expErr)
	res, err := s.uc.TransferMoney(sender.ID, sender.ID, sender.Amount)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteBalanceUsecase) TestTransferMoney_SenderHaveNotMoney() {
	sender := test_data.TestBalance(s.T())

	moneyToTransfer := sender.Amount * 2

	expErr := NotEnoughMoney

	s.MockBalanceRepository.EXPECT().
		FindUserByID(sender.ID).
		Return(sender, nil)

	res, err := s.uc.TransferMoney(sender.ID, sender.ID, moneyToTransfer)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteBalanceUsecase) TestTransferMoney_ReceiverNotFound() {
	sender := test_data.TestBalance(s.T())
	receiver := test_data.TestBalance(s.T())
	receiver.ID = 100

	moneyToTransfer := sender.Amount
	expErr := balance_repository.NotFound

	s.MockBalanceRepository.EXPECT().
		FindUserByID(sender.ID).
		Return(sender, nil)
	s.MockBalanceRepository.EXPECT().
		FindUserByID(receiver.ID).
		Return(nil, expErr)

	res, err := s.uc.TransferMoney(sender.ID, receiver.ID, moneyToTransfer)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteBalanceUsecase) TestTransferMoney_TransferError() {
	sender := test_data.TestBalance(s.T())
	receiver := test_data.TestBalance(s.T())
	receiver.ID = 100

	moneyToTransfer := sender.Amount
	expErr := balance_repository.NewDBError(balance_repository.DefaultErrDB)

	s.MockBalanceRepository.EXPECT().
		FindUserByID(sender.ID).
		Return(sender, nil)
	s.MockBalanceRepository.EXPECT().
		FindUserByID(receiver.ID).
		Return(receiver, nil)

	s.MockBalanceRepository.EXPECT().
		CreateTransfer(sender.ID, receiver.ID, moneyToTransfer).
		Return(expErr)

	res, err := s.uc.TransferMoney(sender.ID, receiver.ID, moneyToTransfer)
	assert.Nil(s.T(), res)
	assert.Equal(s.T(), expErr, err)
}
func (s *SuiteBalanceUsecase) TestTransferMoney_OK() {
	sender := test_data.TestBalance(s.T())
	receiver := test_data.TestBalance(s.T())
	receiver.ID = 100

	moneyToTransfer := sender.Amount
	expRes := &models.TransferMoney{
		SenderID:        sender.ID,
		ReceiverID:      receiver.ID,
		SenderBalance:   sender.Amount - moneyToTransfer,
		ReceiverBalance: receiver.Amount + moneyToTransfer,
	}
	s.MockBalanceRepository.EXPECT().
		FindUserByID(sender.ID).
		Return(sender, nil)
	s.MockBalanceRepository.EXPECT().
		FindUserByID(receiver.ID).
		Return(receiver, nil)

	s.MockBalanceRepository.EXPECT().
		CreateTransfer(sender.ID, receiver.ID, moneyToTransfer).
		Return(nil)

	res, err := s.uc.TransferMoney(sender.ID, receiver.ID, moneyToTransfer)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expRes, res)
}
func TestUsecaseCreator(t *testing.T) {
	suite.Run(t, new(SuiteBalanceUsecase))
}
