package balance_repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

func (s *Suite) InitBD() {
	s.T().Helper()

	var err error
	s.DB, s.Mock, err = sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
}
