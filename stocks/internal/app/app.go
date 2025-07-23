package app

import (
	"context"
	"fmt"
	"log"
	"stocks/internal/app/server"
	"stocks/internal/config"
	"stocks/internal/kafka"
	"stocks/pkg/connection"
	"stocks/pkg/constants"
	"strings"
	"time"
)

func NewStockServiceApp() error {
	// load environment variables.
	if err := config.LoadEnv(".env"); err != nil {
		log.Printf("app.config.LoadEnv: %v+\n", err.Error())
		return fmt.Errorf("failed to load env file: %w", err)
	}

	cfg, err := config.NewStockServiceConfig()
	if err != nil {
		log.Printf("app.config.NewStockServiceConfig: %+v\n", err.Error())
		return fmt.Errorf("failed to initialize NewStockServiceConfig: %w", err)
	}

	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.DBCtxTimeOut*time.Second)
	defer cancel()

	psqlDB, err := connection.NewDB(ctxTimeOut, cfg.DbConfig())
	if err != nil {
		log.Printf("app.connection.NewDB: %+v\n", err.Error())
		return fmt.Errorf("failed to create a new connection DB: %w", err)
	}

	defer func() {
		psqlDB.Close()
		log.Println("postgres connection successfully completed")
	}()

	kafkaProducer, err := kafka.NewStocksServiceEventProducer(strings.Split(cfg.GetKafkaBrokers(), ","))
	if err != nil {
		log.Printf("failed to create stocks service kafka producer: %v\n", err)
	}
	defer kafkaProducer.Close()

	srv := server.NewServer(cfg, psqlDB, kafkaProducer)
	if err := srv.RunServer(); err != nil {
		log.Printf("%+v\n", err.Error())
		return fmt.Errorf("failed to run http server: %w", err)
	}

	return nil
}
