package balance_handler

import (
	"avito-intern/internal/app/balance"
	"avito-intern/internal/app/balance/delivery"
	models2 "avito-intern/internal/app/balance/delivery/models"
	"avito-intern/internal/app/balance/models"
	"avito-intern/internal/app/middlewares"
	"avito-intern/internal/pkg/handler"
	"avito-intern/internal/pkg/utilits"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
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
			utilits.Responder{
				LogObject: utilits.NewLogObject(logger),
			},
		},
	}

	h.router.HandleFunc("/balance/{user_id:[0-9]+}", h.GetBalanceHandler).Methods(http.MethodGet)
	h.router.HandleFunc("/balance/{user_id:[0-9]+}", h.UpdateBalanceHandler).Methods(http.MethodPost)
	h.router.HandleFunc("/transfer", h.TransferMoneyHandler).Methods(http.MethodPost)
	h.router.HandleFunc("/transaction", h.TransactionHandler).Methods(http.MethodPost)

	utilitiesMiddleware := middlewares.NewUtilitiesMiddleware(h.logger)
	h.router.Use(utilitiesMiddleware.UpgradeLogger)
	h.router.Use(utilitiesMiddleware.CheckPanic)

	return h
}

func (h *BalanceHandler) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	res, ok := h.GetInt64FromParam(w, r, "user_id")
	if !ok {
		return
	}
	amount, err := h.usecase.GetBalance(res)
	if err != nil {
		h.UsecaseError(w, r, err, balance.CodeByErrorGetBalance)
		return
	}

	h.Log(r).Debugf("GET_BALANCE_HANDLER: get balance %v user_id = %v")

	h.Respond(w, r, http.StatusOK, models2.RespondBalance{Balance: amount})
}
func (h *BalanceHandler) TransferMoneyHandler(w http.ResponseWriter, r *http.Request) {

}
func (h *BalanceHandler) TransactionHandler(w http.ResponseWriter, r *http.Request) {

}
func (h *BalanceHandler) UpdateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	req := &models2.RequestUpdateBalance{}
	err := h.GetRequestBody(w, r, req, *bluemonday.UGCPolicy())
	if err != nil {
		h.Log(r).Warnf("can not decode body %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, delivery.InvalidBody)
		return
	}
	if req.Type != models.ADD_BALANCE && req.Type != models.DIFF_BALANCE {
		h.Log(r).Warnf("invalid body - incorrect Type of updateBalance operation %s", err)
		h.Error(w, r, http.StatusUnprocessableEntity, delivery.InvalidBody)
		return
	}
	userID, ok := h.GetInt64FromParam(w, r, "user_id")
	if !ok {
		return
	}
	newBalance, err := h.usecase.UpdateBalance(userID, req.Amount, int(req.Type))
	if err != nil {
		h.UsecaseError(w, r, err, balance.CodeByErrorGetBalance)
		return
	}
	h.Respond(w, r, http.StatusOK, models2.RespondBalance{UserID: userID, Balance: newBalance})
}

func (h *BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
