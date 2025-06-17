package app

import (
	"cart/internal/app/server"
	"cart/internal/config"
	"cart/pkg/connection"
	"cart/pkg/constants"
	"context"
	"fmt"
	"log"
	"time"
)

func NewCartServiceApp() error {
	if err := config.LoadEnv(".env"); err != nil {
		log.Printf("app.config.LoadEnv: %v+\n", err.Error())
		return fmt.Errorf("failed to load env file: %w", err)
	}

	cfg, err := config.NewCartServiceConfig()
	if err != nil {
		return fmt.Errorf("failed to initialize NewCartServiceConfig: %w", err)
	}

	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.DBCtxTimeOut*time.Second)
	defer cancel()

	psqlDB, err := connection.NewDB(ctxTimeOut, cfg.DbConfig())
	if err != nil {
		return fmt.Errorf("failed to create a new connection DB: %w", err)
	}

	defer func() {
		psqlDB.Close()
		log.Println("postgres connection successfully completed")
	}()

	srv := server.NewServer(cfg, psqlDB)
	if err := srv.RunHTTPServer(); err != nil {
		log.Printf("%+v\n", err.Error())
		return fmt.Errorf("can not start http server: %w", err)
	}

	return nil
}
