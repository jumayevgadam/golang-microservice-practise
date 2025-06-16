package server

import (
	"net/http"
	v1 "stocks/internal/controller/http/v1"
	"stocks/internal/repository/postgres"
	"stocks/internal/usecase/stocks"
)

func (s *Server) setupRoutes() *http.ServeMux {
	// initialize repository.
	skuRepo := postgres.NewSKURepository(s.psqlDB)
	stockRepo := postgres.NewStockServiceRepository(s.psqlDB)

	// initialize usecase.
	stockUC := stocks.NewStockServiceUseCase(skuRepo, stockRepo)

	// initialize handlers.
	handlers := &v1.Handlers{
		StockServiceHandler: v1.NewstockServiceController(stockUC),
	}

	return v1.MapRoutes(handlers)
}
