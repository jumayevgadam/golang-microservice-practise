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
	zapLogger "stocks/pkg/log/zap"
	"stocks/pkg/tracer"
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

	logger, cleanup, err := zapLogger.NewLogger(cfg.ExternalServices.ObservalityConfig.LogStashHost)
	if err != nil {
		log.Printf("error initializing logger: %v", err)
	}

	defer cleanup()

	tp, err := tracer.InitTracer(cfg.Server.ServiceName)
	if err != nil {
		return fmt.Errorf("error initializing tracer")
	}

	defer func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			logger.Error("error shutting down tracer")
		}
	}()

	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.DBCtxTimeOut*time.Second)
	defer cancel()

	psqlDB, err := connection.NewDB(ctxTimeOut, cfg.DbConfig())
	if err != nil {
		logger.Errorf("app.connection.NewDB: %+v\n", err.Error())
		return fmt.Errorf("failed to create a new connection DB: %w", err)
	}

	defer func() {
		psqlDB.Close()
		logger.Info("postgres connection successfully completed")
	}()

	kafkaProducer, err := kafka.NewStocksServiceEventProducer(strings.Split(cfg.GetKafkaBrokers(), ","))
	if err != nil {
		logger.Errorf("failed to create stocks service kafka producer: %v\n", err)
	}
	defer kafkaProducer.Close()

	srv := server.NewServer(cfg, psqlDB, kafkaProducer, logger)
	if err := srv.RunServer(); err != nil {
		logger.Errorf("%+v\n", err.Error())
		return fmt.Errorf("failed to run http server: %w", err)
	}

	return nil
}
