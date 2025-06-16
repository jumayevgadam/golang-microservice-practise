package app

import (
	"context"
	"log"
	"stocks/internal/app/server"
	"stocks/internal/config"
	"stocks/pkg/connection"
	"stocks/pkg/constants"
	"time"
)

func NewStockServiceApp() error {
	// load environment variables.
	if err := config.LoadEnv(".env"); err != nil {
		log.Printf("app.config.LoadEnv: %v+\n", err.Error())
	}

	cfg, err := config.NewStockServiceConfig()
	if err != nil {
		log.Printf("app.config.NewStockServiceConfig: %+v\n", err.Error())
	}

	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.DBCtxTimeOut*time.Second)
	defer cancel()

	psqlDB, err := connection.NewDB(ctxTimeOut, cfg.DbConfig())
	if err != nil {
		log.Printf("app.connection.NewDB: %+v\n", err.Error())
	}

	defer func() {
		psqlDB.Close()
		log.Println("postgres connection successfully completed")
	}()

	srv := server.NewServer(cfg, psqlDB)
	if err := srv.RunHTTPServer(); err != nil {
		log.Printf("%+v\n", err.Error())
	}

	return nil
}
