package v1

import "net/http"

type Handlers struct {
	StockServiceHandler *stockServiceController
}

func MapRoutes(h *Handlers) *http.ServeMux {
	mux := http.NewServeMux()
	// stock service endpoints.
	mux.HandleFunc("POST /stocks/item/add", h.StockServiceHandler.AddStockItem)
	mux.HandleFunc("POST /stocks/item/delete", h.StockServiceHandler.DeleteStockItem)
	mux.HandleFunc("POST /stocks/item/get", h.StockServiceHandler.GetStockItemBySKU)
	mux.HandleFunc("POST /stocks/list/location", h.StockServiceHandler.ListStockItems)

	return mux
}
