package server

import (
	grpcV1 "cart/internal/controller/grpc/v1"
	httpV1 "cart/internal/controller/http/v1"
	"cart/internal/repository/postgres"
	"cart/internal/service/stockms"
	"cart/internal/usecase/carts"
	pb "cart/pkg/api/cart"
	"net/http"
)

func (s *Server) setupRoutes() *http.ServeMux {
	// repos.
	cartRepo := postgres.NewCartItemRepository(s.psqlDB)

	// services.
	stockService, _ := stockms.NewGRPCStockService(s.cfg.StockServiceGRPCAddress())

	// usecases.
	cartUseCase := carts.NewCartServiceUseCase(stockService, cartRepo, s.kafkaProducer)

	// handlers.
	handlers := &httpV1.Handlers{
		CartServiceHandler: httpV1.NewCartServiceController(cartUseCase),
	}

	return httpV1.MapRoutes(handlers)
}

func (s *Server) registerGRPCServices() {
	// repos.
	cartRepo := postgres.NewCartItemRepository(s.psqlDB)

	// services. ignored error in this place...
	stockService, _ := stockms.NewGRPCStockService(s.cfg.StockServiceGRPCAddress())

	// usecases.
	cartUseCase := carts.NewCartServiceUseCase(stockService, cartRepo, s.kafkaProducer)

	cartGRPCHandler := grpcV1.NewCartGRPCHandler(cartUseCase)

	pb.RegisterCartServiceServer(s.grpcServer, cartGRPCHandler)
}
