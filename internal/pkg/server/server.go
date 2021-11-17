package server

import (
	"avito-intern/configs"
	"avito-intern/internal/app/balance/balance_repository"
	"avito-intern/internal/app/balance/balance_usecase"
	balance_handler "avito-intern/internal/app/balance/delivery/http"
	transaction_handler "avito-intern/internal/app/transaction/delivery/http"
	"avito-intern/internal/app/transaction/transaction_repository"
	"avito-intern/internal/app/transaction/transaction_usecase"
	"avito-intern/internal/pkg/handler"
	"avito-intern/internal/pkg/utilits"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Server struct {
	handler     handler.Handler
	config      *configs.Config
	logger      *logrus.Logger
	connections utilits.ExpectedConnections
}

func NewServer(config *configs.Config, connections utilits.ExpectedConnections, logger *logrus.Logger) *Server {
	return &Server{
		config:      config,
		connections: connections,
		logger:      logger,
	}
}
func (s *Server) checkConnection() error {
	if err := s.connections.SqlConnection.Ping(); err != nil {
		return fmt.Errorf("Can't check connection to sql with error %v ", err)
	}
	s.logger.Info("Success check connection to sql db")
	return nil
}

func (s *Server) Start() error {
	if err := s.checkConnection(); err != nil {
		return err
	}

	router := mux.NewRouter()
	routerApi := router.PathPrefix("/api/v1/").Subrouter()
	routerApi.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	balanceRepository := balance_repository.NewBalanceRepository(s.connections.SqlConnection)
	balanceUsecase := balance_usecase.NewBalanceUsecase(balanceRepository)

	transactionRepository := transaction_repository.NewTransactionRepository(s.connections.SqlConnection)
	transactionUsecase := transaction_usecase.NewTransactionUsecase(transactionRepository)

	hBalance := balance_handler.NewBalanceHandler(routerApi, s.logger, balanceUsecase, s.config.CurrencyAPI)
	hTransactions := transaction_handler.NewTransactionHandler(routerApi, s.logger, transactionUsecase)

	hBalance.BaseLog().Info("initialize Balance Handler success")
	hTransactions.BaseLog().Info("initialize Transaction handler success")

	s.logger.Info("Server start")
	return http.ListenAndServe(s.config.BindAddr, routerApi)
}
