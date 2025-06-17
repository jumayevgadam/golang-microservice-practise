package v1

import "net/http"

type Handlers struct {
	CartServiceHandler *cartServiceController
}

func MapRoutes(h *Handlers) *http.ServeMux {
	mux := http.NewServeMux()
	// cart service endpoints.
	mux.HandleFunc("POST /cart/item/add", h.CartServiceHandler.AddCartItem)
	mux.HandleFunc("POST /cart/item/delete", h.CartServiceHandler.DeleteCartItem)
	mux.HandleFunc("POST /cart/clear", h.CartServiceHandler.ClearCartItems)
	mux.HandleFunc("POST /cart/list", h.CartServiceHandler.ListCartItems)

	return mux
}
