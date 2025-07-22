package v1

import (
	"errors"
	"net/http"
	"stocks/internal/domain"
	"stocks/internal/usecase"
	"stocks/pkg/httphelper"
)

type stockServiceController struct {
	stockItemUC usecase.StockServiceUseCase
}

func NewstockServiceController(stockItemUC usecase.StockServiceUseCase) *stockServiceController {
	return &stockServiceController{stockItemUC: stockItemUC}
}

func (c *stockServiceController) AddStockItem(w http.ResponseWriter, r *http.Request) {
	var stockItemReq CreateStockItemRequest
	if err := httphelper.RequestValidate(r, &stockItemReq); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())

		return
	}
	// call usecase.
	err := c.stockItemUC.AddStockItem(r.Context(), stockItemReq.ToDomain())
	if err != nil {
		if errors.Is(err, domain.ErrSKUNotFound) {
			httphelper.ErrorResponse(w, http.StatusNotFound, "SKU not found")
			return
		}

		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}
	// give response
	httphelper.Respond(w, http.StatusOK, map[string]string{"message": "stock item added or updated successfully"})
}

func (c *stockServiceController) DeleteStockItem(w http.ResponseWriter, r *http.Request) {
	var deleteStockItemReq DeleteStockItemRequest
	if err := httphelper.RequestValidate(r, &deleteStockItemReq); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err := c.stockItemUC.DeleteStockItem(r.Context(), domain.UserID(deleteStockItemReq.UserID), domain.SKUID(deleteStockItemReq.SkuID))
	if err != nil {
		if errors.Is(err, domain.ErrStockItemNotFound) {
			httphelper.ErrorResponse(w, http.StatusNotFound, "stock item not found")
			return
		}

		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	httphelper.Respond(w, http.StatusOK, map[string]string{"message": "stock item deleted successfully"})
}

func (c *stockServiceController) GetStockItemBySKU(w http.ResponseWriter, r *http.Request) {
	var getStocItemBySKUReq GetStockItemRequest
	if err := httphelper.RequestValidate(r, &getStocItemBySKUReq); err != nil {
		httphelper.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	stockItem, err := c.stockItemUC.GetStockItemBySKU(r.Context(), domain.SKUID(getStocItemBySKUReq.SkuID))
	if err != nil {
		if errors.Is(err, domain.ErrStockItemNotFound) {
			httphelper.ErrorResponse(w, http.StatusNotFound, "stock_item not found")
			return
		}

		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	resp := GetStockItemResponse{
		SkuID:    uint32(stockItem.Sku.ID),
		Name:     stockItem.Sku.Name,
		Type:     stockItem.Sku.Type,
		Count:    stockItem.Count,
		Price:    stockItem.Price,
		Location: stockItem.Location,
	}

	httphelper.Respond(w, http.StatusOK, resp)
}

func (c *stockServiceController) ListStockItems(w http.ResponseWriter, r *http.Request) {
	var filteredRequest FilterRequest
	if err := httphelper.RequestValidate(r, &filteredRequest); err != nil {
		httphelper.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	listStockItemsResult, err := c.stockItemUC.ListStockItems(r.Context(), filteredRequest.ToDomain())
	if err != nil {
		httphelper.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := ToListResponse(listStockItemsResult)
	httphelper.Respond(w, http.StatusOK, response)
}
