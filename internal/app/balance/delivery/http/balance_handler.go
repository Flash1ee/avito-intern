package balance_handler

import (
	"avito-intern/internal/app/balance"
	request_response_models "avito-intern/internal/app/balance/delivery/models"
	"avito-intern/internal/app/middlewares"
	"avito-intern/internal/pkg/handler"
	"avito-intern/internal/pkg/utilits"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type BalanceHandler struct {
	router  *mux.Router
	logger  *logrus.Logger
	usecase balance.Usecase
	handler.HelpHandlers
}

func NewBalanceHandler(router *mux.Router, logger *logrus.Logger, uc balance.Usecase) *BalanceHandler {
	h := &BalanceHandler{
		router:  router,
		logger:  logger,
		usecase: uc,
		HelpHandlers: handler.HelpHandlers{
			Responder: utilits.Responder{
				LogObject: utilits.NewLogObject(logger),
			},
		},
	}

	h.router.HandleFunc("/balance/{user_id}", h.GetBalanceHandler).Methods(http.MethodGet)
	h.router.HandleFunc("/balance/{user_id}", h.UpdateBalanceHandler).Methods(http.MethodPost)
	h.router.HandleFunc("/transfer", h.TransferMoneyHandler).Methods(http.MethodPost)
	//h.router.HandleFunc("/transaction", h.TransactionHandler).Methods(http.MethodPost)

	utilitiesMiddleware := middlewares.NewUtilitiesMiddleware(h.logger)
	h.router.Use(utilitiesMiddleware.UpgradeLogger)
	h.router.Use(utilitiesMiddleware.CheckPanic)

	return h
}

// GetBalanceHandler
// @Summary get user balance
// @Description get user balance with id from query
// @Produce json
// @Param user_id path int true "user_id in balanceApp"
// @Success 200 {object} models.ResponseBalance
// @Failure 400 {object} models.ErrResponse "invalid query param"
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 500 {object} models.ErrResponse "internal error"
// @Router /balance/{:user_id} [GET]
func (h *BalanceHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.GetInt64FromParam(w, r, "user_id")
	if !ok {
		return
	}
	amount, err := h.usecase.GetBalance(userID)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorGetBalance)
		return
	}

	h.Log(r).Debugf("GET_BALANCE_HANDLER: get balance %v user_id = %v", amount, userID)

	h.Respond(w, r, http.StatusOK, request_response_models.ResponseBalance{Balance: amount})
}

// TransferMoneyHandler
// @Summary transfer money
// @Description money transfer between users
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseTransfer "Successfully transaction"
// @Failure 404 {object} models.ErrResponse "user with one of id in request body not found"
// @Failure 422 {object} models.ErrResponse "invalid body in request || not enough money for transfer"
// @Failure 500 {object} models.ErrResponse "internal error"
// @Router /transfer [POST]
func (h *BalanceHandler) TransferMoneyHandler(w http.ResponseWriter, r *http.Request) {
	req := &request_response_models.RequestTransfer{}
	err := h.GetRequestBody(w, r, req)
	if err != nil {
		h.Log(r).Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler.InvalidBody)
		return
	}
	if err = req.Validate(); err != nil {
		h.Log(r).Warnf("invalid RequestTransferMoney body err: %v body: %v", err, req)
		h.Error(w, r, http.StatusUnprocessableEntity, handler.InvalidBody)
		return
	}
	res, err := h.usecase.TransferMoney(req.SenderID, req.ReceiverID, req.Amount)
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorTransferHandler)
		return
	}
	h.Respond(w, r, http.StatusOK, request_response_models.ResponseTransfer{
		SenderID:        res.SenderID,
		SenderBalance:   res.SenderBalance,
		ReceiverID:      res.ReceiverID,
		ReceiverBalance: res.ReceiverBalance})
}

// UpdateBalanceHandler
// @Summary update user balance
// @Description update user balance
// @Produce json
// @Param data body models.RequestUpdateBalance true "data for update balance, operation = 0 - writeOff money, operation = 1 - addMoney"
// @Param user_id path int true "user_id in balanceApp"
// @Success 200 {object} models.ResponseBalance
// @Failure 400 {object} models.ErrResponse "invalid query param"
// @Failure 404 {object} models.ErrResponse "user with this id not found"
// @Failure 422 {object} models.ErrResponse "invalid body in request || not enough money for transfer"
// @Failure 500 {object} models.ErrResponse "internal error"
// @Router /balance/{:user_id} [POST]
func (h *BalanceHandler) UpdateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	req := &request_response_models.RequestUpdateBalance{}
	err := h.GetRequestBody(w, r, req)
	if err != nil {
		h.Log(r).Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, handler.InvalidBody)
		return
	}
	if err = req.Validate(); err != nil {
		h.Log(r).Warnf("invalid RequestUpdateBalance body err: %v body: %v", err, req)
		h.Error(w, r, http.StatusUnprocessableEntity, handler.InvalidBody)
		return
	}
	userID, ok := h.GetInt64FromParam(w, r, "user_id")
	if !ok {
		return
	}
	newBalance, err := h.usecase.UpdateBalance(userID, req.Amount, int(req.Type))
	if err != nil {
		h.UsecaseError(w, r, err, CodeByErrorBalanceHandler)
		return
	}
	h.Respond(w, r, http.StatusOK, request_response_models.ResponseBalance{UserID: userID, Balance: newBalance})
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
