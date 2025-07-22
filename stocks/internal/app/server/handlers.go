package server

import (
	"net/http"
	grpcV1 "stocks/internal/controller/grpc/v1"
	httpV1 "stocks/internal/controller/http/v1"
	"stocks/internal/repository/postgres"
	stockUC "stocks/internal/usecase/stocks"
	pb "stocks/pkg/api/stocks"
)

func (s *Server) setupRoutes() *http.ServeMux {
	// initialize repository.
	skuRepo := postgres.NewSKURepository(s.psqlDB)
	stockRepo := postgres.NewStockServiceRepository(s.psqlDB)

	// initialize usecase.
	stockUC := stockUC.NewStockServiceUseCase(skuRepo, stockRepo, s.kafkaProducer)

	// initialize handlers.
	handlers := &httpV1.Handlers{
		StockServiceHandler: httpV1.NewstockServiceController(stockUC),
	}

	return httpV1.MapRoutes(handlers)
}

func (s *Server) registerGRPCServices() {
	// initialize repository.
	skuRepo := postgres.NewSKURepository(s.psqlDB)
	stockRepo := postgres.NewStockServiceRepository(s.psqlDB)

	// initialize usecase.
	stockUC := stockUC.NewStockServiceUseCase(skuRepo, stockRepo, s.kafkaProducer)

	stockGRPCHandler := grpcV1.NewStockGRPCHandler(stockUC)

	pb.RegisterStocksServiceServer(s.grpcServer, stockGRPCHandler)
}
