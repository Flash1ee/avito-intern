package transaction_handler

import (
	"avito-intern/internal/app/middlewares"
	"avito-intern/internal/app/transaction"
	request_response_models "avito-intern/internal/app/transaction/delivery/models"
	"avito-intern/internal/app/transaction/models"
	"avito-intern/internal/app/transaction/transaction_repository"
	"avito-intern/internal/pkg/handler"
	"avito-intern/internal/pkg/utilits"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	router  *mux.Router
	logger  *logrus.Logger
	usecase transaction.Usecase
	handler.HelpHandlers
}

func NewTransactionHandler(router *mux.Router, logger *logrus.Logger,
	uc transaction.Usecase) *TransactionHandler {
	h := &TransactionHandler{
		router:  router,
		logger:  logger,
		usecase: uc,
		HelpHandlers: handler.HelpHandlers{
			Responder: utilits.Responder{
				LogObject: utilits.NewLogObject(logger),
			},
		},
	}

	h.router.HandleFunc("/transaction/{user_id}", h.GetTransactionHandler).Methods(http.MethodGet)

	utilitiesMiddleware := middlewares.NewUtilitiesMiddleware(h.logger)
	h.router.Use(utilitiesMiddleware.UpgradeLogger)
	h.router.Use(utilitiesMiddleware.CheckPanic)

	return h
}

// GetTransactionHandler
// @Summary get user transactions
// @Produce json
// @Param user_id path int true "user_id in balanceApp"
// @Param count query int true "count of transactions to response"
// @Param page query int true "page of transactions"
// @Param sort query string false "sort for transactions: date(created date) or sum(amount of transaction)"
// @Param direction query string false "sort transactions: asc/desc; usage only with sort param!"
// @Success 200 {object} models.ResponseTransactions
// @Success 204 {object} models.NotFoundResponse "user transactions not found"
// @Failure 400 {object} models.ErrResponse "invalid query param | usage direction without sort"
// @Failure 500 {object} models.ErrResponse "internal error"
// @Router /transaction/{:user_id} [GET]
func (h *TransactionHandler) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := h.GetInt64FromParam(w, r, "user_id")
	if !ok {
		return
	}
	if userID < 0 {
		h.Log(r).Infof("invalid user_id %v", userID)
		h.Error(w, r, http.StatusBadRequest, handler.InvalidParameters)
		return
	}

	var err error
	paginator := &models.Paginator{}

	page := r.URL.Query().Get(PAGE)
	paginator.Page, err = strconv.Atoi(page)
	if err != nil || paginator.Page < 1 {
		h.Error(w, r, http.StatusBadRequest, InvalidQueryPageParam)
		return
	}

	count := r.URL.Query().Get(COUNT)

	paginator.Count, err = strconv.Atoi(count)
	if err != nil || paginator.Count < 1 {
		h.Error(w, r, http.StatusBadRequest, InvalidQueryCountParam)
		return
	}
	sortField := r.URL.Query().Get(SORT)
	switch sortField {
	case DATE:
		paginator.SortField = models.DATE
	case SUM:
		paginator.SortField = models.SUM
	default:
		if sortField != "" {
			h.Error(w, r, http.StatusBadRequest, InvalidQuerySortParam)
			return
		}
		paginator.SortField = models.NO_ORDER
	}

	directionSort := r.URL.Query().Get(DIRECTION)
	switch directionSort {
	case ASC:
		paginator.SortDirection = models.ASC
	case DESC:
		paginator.SortDirection = models.DESC
	default:
		if directionSort != "" {
			h.Error(w, r, http.StatusBadRequest, InvalidQueryDirectionParam)
			return
		}
		paginator.SortDirection = models.NO_DIRECTION
	}
	if paginator.SortField == models.NO_ORDER &&
		paginator.SortDirection != models.NO_DIRECTION {
		h.Error(w, r, http.StatusBadRequest, DirectionMustUsageWithSort)
		return

	}

	userTransactions, err := h.usecase.GetTransactions(userID, paginator)
	if err != nil {
		if err == transaction_repository.NotFound {
			h.Respond(w, r, http.StatusNoContent, request_response_models.NotFoundResponse{
				Ok: TransactionsNotFound.Error(),
			})
		} else {
			h.UsecaseError(w, r, err, codeByErrorGetTransactions)
		}
		return
	}
	res := request_response_models.ToResponseTransactions(userTransactions)

	h.Respond(w, r, http.StatusOK, res)
}
func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}
