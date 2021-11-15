package server

import (
	"avito-intern/configs"
	"avito-intern/internal/app/balance/balance_repository"
	"avito-intern/internal/app/balance/balance_usecase"
	balance_handler "avito-intern/internal/app/balance/delivery/http"
	"avito-intern/internal/pkg/handler"
	"avito-intern/internal/pkg/utilits"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	balanceRepository := balance_repository.NewBalanceRepository(s.connections.SqlConnection)
	balanceUsecase := balance_usecase.NewBalanceUsecase(balanceRepository)

	h := balance_handler.NewBalanceHandler(router, s.logger, balanceUsecase)
	s.logger.Info("Server start")
	return http.ListenAndServe(s.config.BindAddr, h)
}
