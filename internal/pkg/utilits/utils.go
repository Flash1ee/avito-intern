package utilits

import (
	"avito-intern/configs"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type ExpectedConnections struct {
	SqlConnection *sql.DB
	PathFiles     string
}

func NewLogger(config *configs.Config) (log *logrus.Logger, closeResource func() error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Fatal(err)
	}

	logger := logrus.New()
	currentTime := time.Now().In(time.UTC)

	formatted := config.LogAddr + fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), currentTime.Minute(), currentTime.Second()) + ".log"

	f, err := os.OpenFile(formatted, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("error opening file: %v", err)
	}

	logger.SetOutput(f)
	logger.Writer()
	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger, f.Close
}

func NewPostgresConnection(config *configs.RepositoryConnections) (db *sql.DB, closeResource func() error) {
	db, err := sql.Open("postgres", config.DataBaseUrl)
	if err != nil {
		logrus.Fatal(err)
	}

	return db, db.Close
}
