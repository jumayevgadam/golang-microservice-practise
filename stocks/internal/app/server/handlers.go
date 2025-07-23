package server

import (
	grpcV1 "stocks/internal/controller/grpc/v1"
	"stocks/internal/repository/postgres"
	stockUC "stocks/internal/usecase/stocks"
	pb "stocks/pkg/api/stocks"
)

func (s *Server) registerGRPCServices() {
	// initialize repository.
	skuRepo := postgres.NewSKURepository(s.psqlDB)
	stockRepo := postgres.NewStockServiceRepository(s.psqlDB)

	// initialize usecase.
	stockUC := stockUC.NewStockServiceUseCase(skuRepo, stockRepo, s.kafkaProducer)

	stockGRPCHandler := grpcV1.NewStockGRPCHandler(stockUC)

	pb.RegisterStocksServiceServer(s.grpcServer, stockGRPCHandler)
}
