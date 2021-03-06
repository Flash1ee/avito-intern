package main

import (
	"avito-intern/configs"
	_ "avito-intern/docs"
	"avito-intern/internal/pkg/server"
	"avito-intern/internal/pkg/utilits"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

// @title BalanceApp
// @version 1.0
// @description Server for Balance application.

// @BasePath /api/v1

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	flag.Parse()
	config := configs.NewConfig()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
	}

	logger, closeResource := utilits.NewLogger(config)
	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeResource, logger)

	db, closeDbResource := utilits.NewPostgresConnection(&config.ServerRepository)

	defer func(closer func() error, log *logrus.Logger) {
		err := closer()
		if err != nil {
			log.Fatal(err)
		}
	}(closeDbResource, logger)

	serv := server.NewServer(config,
		utilits.ExpectedConnections{
			SqlConnection: db,
		},
		logger,
	)
	if err = serv.Start(); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Server was stopped")

}
