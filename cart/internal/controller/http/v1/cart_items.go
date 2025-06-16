package v1

import (
	"cart/internal/domain"
	"cart/internal/usecase"
	"errors"
	"net/http"
	"stocks/pkg/httphelper"
)

type cartServiceController struct {
	cartServiceUC usecase.CartItemUseCase
}

func NewCartServiceController(cartServiceUC usecase.CartItemUseCase) *cartServiceController {
	return &cartServiceController{cartServiceUC: cartServiceUC}
}

func (c *cartServiceController) AddCartItem(w http.ResponseWriter, r *http.Request) {
	var createCartItemRequest CreateCartItemRequest
	if err := httphelper.RequestValidate(r, &createCartItemRequest); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.cartServiceUC.AddCartItem(r.Context(), createCartItemRequest.ToDomain())
	if err != nil {
		if errors.Is(err, domain.ErrInSufficientStockCount) {
			httphelper.ErrorResponse(w, http.StatusPreconditionFailed, "insufficient stock count")
			return
		}

		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	httphelper.Respond(w, http.StatusOK, map[string]string{"message": "cart item successfully added"})
}

func (c *cartServiceController) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	var deleteItemReq DeleteCartItemRequest
	if err := httphelper.RequestValidate(r, &deleteItemReq); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.cartServiceUC.DeleteCartItem(r.Context(), domain.UserID(deleteItemReq.UserID), domain.SkuID(deleteItemReq.SkuID))
	if err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			httphelper.ErrorResponse(w, http.StatusNotFound, "cart item not found")
			return
		}

		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	httphelper.Respond(w, http.StatusOK, map[string]string{"message": "cart item successfully deleted"})
}

func (c *cartServiceController) ClearCartItems(w http.ResponseWriter, r *http.Request) {
	var clearCartItemRequest ClearCartItemRequest
	if err := httphelper.RequestValidate(r, &clearCartItemRequest); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.cartServiceUC.ClearCartItems(r.Context(), domain.UserID(clearCartItemRequest.UserID))
	if err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			httphelper.ErrorResponse(w, http.StatusNotFound, "cart_items not found")
			return
		}

		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	httphelper.Respond(w, http.StatusOK, map[string]string{"message": "cart items successfully deleted"})
}

func (c *cartServiceController) ListCartItems(w http.ResponseWriter, r *http.Request) {
	var listCartItemsRequest ListCartItemsRequest
	if err := httphelper.RequestValidate(r, &listCartItemsRequest); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	listCartItems, err := c.cartServiceUC.ListCartItems(r.Context(), domain.UserID(listCartItemsRequest.UserID))
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ToListResponse(listCartItems)
	httphelper.Respond(w, http.StatusOK, response)
}
