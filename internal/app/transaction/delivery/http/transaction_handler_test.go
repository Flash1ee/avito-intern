package transaction_handler

import (
	test_data "avito-intern/internal/app/balance/testing"
	"avito-intern/internal/app/transaction/delivery/models"
	models_constants "avito-intern/internal/app/transaction/models"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type SuiteTransactionHandler struct {
	SuiteHandler
	handler *TransactionHandler
}

func (s *SuiteTransactionHandler) SetupSuite() {
	s.SuiteHandler.SetupSuite()
	s.handler = NewTransactionHandler(mux.NewRouter(), s.Logger, s.MockTransactionUsecase)
}

func (s *SuiteTransactionHandler) TestTransactionHandler_GetTransactionHandler_OK() {
	tb := &TestTable{
		Name:              "Correct work GetTransactionHandler",
		Data:              struct{}{},
		ExpectedMockTimes: 1,
		ExpectedCode:      http.StatusOK,
	}
	userID := int64(1)
	paginator := test_data.TestPaginator(s.T())
	ucTransactions := test_data.TestTransactions(s.T())

	recorder := httptest.NewRecorder()
	path := "/transaction/1?page=1&count=2"
	paginator.Page = 1
	paginator.Count = 2
	paginator.SortDirection = models_constants.NO_DIRECTION
	paginator.SortField = models_constants.NO_ORDER

	s.MockTransactionUsecase.EXPECT().
		GetTransactions(userID, paginator).
		Return(ucTransactions, nil)

	b := bytes.Buffer{}
	reader, _ := http.NewRequest(http.MethodGet, path, &b)
	s.handler.ServeHTTP(recorder, reader)

	assert.Equal(s.T(), tb.ExpectedCode, recorder.Code)
	req := &models.ResponseTransactions{}
	decoder := json.NewDecoder(recorder.Body)
	err := decoder.Decode(req)
	assert.NoError(s.T(), err)

	assert.Equal(s.T(), req.Transactions[0].UserID, ucTransactions[0].UserID)
	assert.Equal(s.T(), req.Transactions[0].Type, ucTransactions[0].Type)
	assert.Equal(s.T(), req.Transactions[0].Description, ucTransactions[0].Description)
	assert.Zero(s.T(), req.Transactions[0].CreatedAt.Sub(ucTransactions[0].CreatedAt))
	assert.Equal(s.T(), req.Transactions[1].UserID, ucTransactions[1].UserID)
	assert.Equal(s.T(), req.Transactions[1].Type, ucTransactions[1].Type)
	assert.Equal(s.T(), req.Transactions[1].ReceiverID, ucTransactions[1].ReceiverID)
	assert.Equal(s.T(), req.Transactions[1].Description, ucTransactions[1].Description)
	assert.Zero(s.T(), req.Transactions[1].CreatedAt.Sub(ucTransactions[1].CreatedAt))

}
func TestBalanceHandler(t *testing.T) {
	suite.Run(t, new(SuiteTransactionHandler))
}
