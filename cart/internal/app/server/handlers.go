package server

import (
	grpcV1 "cart/internal/controller/grpc/v1"
	"cart/internal/repository/postgres"
	"cart/internal/service/stockms"
	"cart/internal/usecase/carts"
	pb "cart/pkg/api/cart"
	"fmt"
)

func (s *Server) registerGRPCServices() error {
	// repos.
	cartRepo := postgres.NewCartItemRepository(s.psqlDB)

	// services.
	stockService, err := stockms.NewGRPCStockService(s.cfg.StockServiceGRPCAddress())
	if err != nil {
		return fmt.Errorf("failed to create new gRPC stock service: %w", err)
	}

	// usecases.
	cartUseCase := carts.NewCartServiceUseCase(stockService, cartRepo, s.kafkaProducer)

	cartGRPCHandler := grpcV1.NewCartGRPCHandler(cartUseCase, s.logger)

	pb.RegisterCartServiceServer(s.grpcServer, cartGRPCHandler)

	return nil
}
