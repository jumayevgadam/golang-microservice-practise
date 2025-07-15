package server

import (
	v1 "cart/internal/controller/http/v1"
	"cart/internal/repository/postgres"
	"cart/internal/service/stockms"
	"cart/internal/usecase/carts"
	"net/http"
)

func (s *Server) setupRoutes() *http.ServeMux {
	// repos.
	cartRepo := postgres.NewCartItemRepository(s.psqlDB)

	// services.
	stockService := stockms.NewHTTPStockService(s.cfg.StockServiceURL())

	// usecases.
	cartUseCase := carts.NewCartServiceUseCase(stockService, cartRepo, s.kafkaProducer)

	// handlers.
	handlers := &v1.Handlers{
		CartServiceHandler: v1.NewCartServiceController(cartUseCase),
	}

	return v1.MapRoutes(handlers)
}
