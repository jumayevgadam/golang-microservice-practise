package app

import (
	"cart/internal/app/server"
	"cart/internal/config"
	"cart/internal/kafka"
	"cart/pkg/connection"
	"cart/pkg/constants"
	zapLogger "cart/pkg/log/zap"
	"cart/pkg/tracer"
	"context"
	"fmt"
	"log"
	"strings"
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

	// init logger.
	logger, cleanup, err := zapLogger.NewLogger()
	if err != nil {
		log.Printf("error initializing logger: %v\n", err.Error())
		return fmt.Errorf("error initializing logger: %w", err)
	}

	defer cleanup()

	// init tracer.
	tp, err := tracer.InitTracer(cfg.Server.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
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
		return fmt.Errorf("failed to create a new connection DB: %w", err)
	}

	defer func() {
		psqlDB.Close()
		log.Println("postgres connection successfully completed")
	}()

	kafkaProducer, err := kafka.NewCartServiceProducer(strings.Split(cfg.GetKafkaBrokers(), ","))
	if err != nil {
		logger.Errorf("failed to initialize cart service kafka producer: %v\n", err.Error())
	}
	defer kafkaProducer.Close()

	srv := server.NewServer(cfg, psqlDB, kafkaProducer, logger)
	if err := srv.RunServer(); err != nil {
		logger.Errorf("%+v\n", err.Error())
		return fmt.Errorf("can not start http server: %w", err)
	}

	return nil
}
