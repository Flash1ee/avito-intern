package balance_handler

import (
	"avito-intern/internal/app"
	"avito-intern/internal/app/balance/balance_repository"
	test_data "avito-intern/internal/app/balance/balance_repository/testing"
	request_response_models "avito-intern/internal/app/balance/delivery/models"
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
func TestBalanceHandler(t *testing.T) {
	suite.Run(t, new(SuiteBalanceHandler))
}
