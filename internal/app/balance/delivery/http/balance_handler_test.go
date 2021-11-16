package balance_handler

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance/balance_repository"
	test_data "avito-intern/internal/app/balance/balance_repository/testing"
	"avito-intern/internal/app/balance/balance_usecase"
	request_response_models "avito-intern/internal/app/balance/delivery/models"
	"avito-intern/internal/app/balance/models"
	"avito-intern/internal/pkg/handler"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type SuiteBalanceHandler struct {
	SuiteHandler
	handler *BalanceHandler
}

func (s *SuiteBalanceHandler) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewBalanceHandler(mux.NewRouter(), s.Logger, s.MockBalanceUsecase)
}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetBalanceHandler_OK() {
	tb := &TestTable{
		Name:              "Correct work GetBalanceHandler",
		Data:              struct{}{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}
	userBalance := test_data.TestBalance(s.T())

	s.MockBalanceUsecase.EXPECT().GetBalance(userBalance.ID).
		Return(userBalance.Amount, nil).Times(tb.ExpectedMockTimes)
	recorder := httptest.NewRecorder()

	userID := int64(1)
	toStr := strconv.Itoa(int(userID))

	path := "/balance/" + toStr
	b := bytes.Buffer{}
	reader, _ := http.NewRequest(http.MethodGet, path, &b)
	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	req := &request_response_models.ResponseBalance{}
	decoder := json.NewDecoder(recorder.Body)
	err := decoder.Decode(req)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), req, &request_response_models.ResponseBalance{
		Balance: userBalance.Amount,
	})
}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetBalanceHandler_InvalidQueryParam() {
	tb := &TestTable{
		Name:              "Invalid param in query",
		Data:              struct{}{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusBadRequest,
	}

	recorder := httptest.NewRecorder()
	path := "/balance/f"
	b := bytes.Buffer{}
	reader, _ := http.NewRequest(http.MethodGet, path, &b)
	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)

}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetBalanceHandler_GetBalanceError() {
	tb := &TestTable{
		Name:              "Usecase error on get balance - balance not found",
		Data:              struct{}{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusNotFound,
	}
	recorder := httptest.NewRecorder()
	userID := int64(1)
	toStr := strconv.Itoa(int(userID))

	path := "/balance/" + toStr
	b := bytes.Buffer{}
	reader, _ := http.NewRequest(http.MethodGet, path, &b)
	s.MockBalanceUsecase.EXPECT().GetBalance(userID).
		Times(tb.ExpectedMockTimes).
		Return(app.InvalidFloat, balance_repository.NotFound)

	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)

}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetBalanceHandler_GetBalanceError_Internal() {
	tb := &TestTable{
		Name:              "Usecase error on get balance - internal error",
		Data:              struct{}{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	recorder := httptest.NewRecorder()
	userID := int64(1)
	toStr := strconv.Itoa(int(userID))

	path := "/balance/" + toStr
	b := bytes.Buffer{}
	reader, _ := http.NewRequest(http.MethodGet, path, &b)
	s.MockBalanceUsecase.EXPECT().GetBalance(userID).
		Times(tb.ExpectedMockTimes).
		Return(app.InvalidFloat, balance_repository.DefaultErrDB)

	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)

}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetTransferMoneyHandler_OK() {
	req := request_response_models.RequestTransfer{
		SenderID:   1,
		ReceiverID: 2,
		Amount:     100.5,
	}
	transferResponse := &models.TransferMoney{
		SenderID:        req.SenderID,
		ReceiverID:      req.ReceiverID,
		SenderBalance:   req.Amount - req.Amount/2,
		ReceiverBalance: req.Amount + req.Amount/2,
	}
	tb := &TestTable{
		Name:              "Successfully work",
		Data:              req,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(tb.Data)
	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/transfer", &b)
	s.MockBalanceUsecase.EXPECT().TransferMoney(req.SenderID, req.ReceiverID, req.Amount).
		Times(tb.ExpectedMockTimes).
		Return(transferResponse, nil)

	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	responseRes := &request_response_models.ResponseTransfer{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(responseRes)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), transferResponse.SenderID, responseRes.SenderID)
	assert.Equal(s.T(), transferResponse.ReceiverID, responseRes.ReceiverID)
	assert.Equal(s.T(), transferResponse.ReceiverBalance, responseRes.ReceiverBalance)
	assert.Equal(s.T(), transferResponse.SenderBalance, responseRes.SenderBalance)

}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetTransferMoneyHandler_UserNotFound() {
	req := request_response_models.RequestTransfer{
		SenderID:   1,
		ReceiverID: 2,
		Amount:     100.5,
	}
	tb := &TestTable{
		Name:              "Sender/receiver not found id DB",
		Data:              req,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusNotFound,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(tb.Data)
	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/transfer", &b)
	s.MockBalanceUsecase.EXPECT().TransferMoney(req.SenderID, req.ReceiverID, req.Amount).
		Times(tb.ExpectedMockTimes).
		Return(nil, balance_repository.NotFound)

	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	responseRes := &request_response_models.ErrResponse{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(responseRes)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), responseRes.Err, handler.UserNotFound.Error())
}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetTransferMoneyHandler_InternalError() {
	req := request_response_models.RequestTransfer{
		SenderID:   1,
		ReceiverID: 2,
		Amount:     100.5,
	}
	tb := &TestTable{
		Name:              "Database Error",
		Data:              req,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusInternalServerError,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(tb.Data)
	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/transfer", &b)
	s.MockBalanceUsecase.EXPECT().TransferMoney(req.SenderID, req.ReceiverID, req.Amount).
		Times(tb.ExpectedMockTimes).
		Return(nil, balance_repository.DefaultErrDB)

	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	responseRes := &request_response_models.ErrResponse{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(responseRes)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), responseRes.Err, handler.BDError.Error())
}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetTransferMoneyHandler_SenderNoHaveMoney() {
	req := request_response_models.RequestTransfer{
		SenderID:   1,
		ReceiverID: 2,
		Amount:     100.5,
	}
	tb := &TestTable{
		Name:              "Sum of transfer more then sender have on balance",
		Data:              req,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(tb.Data)
	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/transfer", &b)
	s.MockBalanceUsecase.EXPECT().TransferMoney(req.SenderID, req.ReceiverID, req.Amount).
		Times(tb.ExpectedMockTimes).
		Return(nil, balance_usecase.NotEnoughMoney)

	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	responseRes := &request_response_models.ErrResponse{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(responseRes)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), responseRes.Err, handler.NotEnoughMoney.Error())
}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetTransferMoneyHandler_InvalidBody() {
	req := request_response_models.RequestTransfer{
		SenderID:   -1,
		ReceiverID: 2,
		Amount:     100.5,
	}
	tb := &TestTable{
		Name:              "Invalid id in body",
		Data:              req,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(tb.Data)
	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/transfer", &b)
	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	responseRes := &request_response_models.ErrResponse{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(responseRes)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), responseRes.Err, handler.InvalidBody.Error())
}
func (s *SuiteBalanceHandler) TestBalanceHandler_GetTransferMoneyHandler_InvalidBodyData() {
	req := models.Balance{}
	tb := &TestTable{
		Name:              "Invalid type of body",
		Data:              req,
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusUnprocessableEntity,
	}
	recorder := httptest.NewRecorder()

	b := bytes.Buffer{}
	err := json.NewEncoder(&b).Encode(tb.Data)
	assert.NoError(s.T(), err)

	reader, _ := http.NewRequest(http.MethodPost, "/transfer", &b)
	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	responseRes := &request_response_models.ErrResponse{}
	decoder := json.NewDecoder(recorder.Body)
	err = decoder.Decode(responseRes)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), responseRes.Err, handler.InvalidBody.Error())
}

func TestBalanceHandler(t *testing.T) {
	suite.Run(t, new(SuiteBalanceHandler))
}
