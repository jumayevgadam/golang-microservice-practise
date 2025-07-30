package app

import (
	"cart/internal/app/server"
	"cart/internal/config"
	"cart/internal/kafka"
	"cart/pkg/connection"
	"cart/pkg/constants"
	zapLogger "cart/pkg/log/zap"
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

	logger, cleanup, err := zapLogger.NewLogger(cfg.Observality.LogStashHost)
	if err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}

	defer cleanup()

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
		log.Printf("failed to initialize cart service kafka producer: %v\n", err.Error())
	}
	defer kafkaProducer.Close()

	srv := server.NewServer(cfg, psqlDB, kafkaProducer, logger)
	if err := srv.RunServer(); err != nil {
		log.Printf("%+v\n", err.Error())
		return fmt.Errorf("can not start http server: %w", err)
	}

	return nil
}
