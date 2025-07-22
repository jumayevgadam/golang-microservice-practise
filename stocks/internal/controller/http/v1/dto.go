package v1

import "stocks/internal/domain"

type CreateStockItemRequest struct {
	UserID   int64  `json:"userID" validate:"required"`
	SkuID    uint32 `json:"skuID" validate:"required"`
	Count    uint16 `json:"count" validate:"required"`
	Price    uint32 `json:"price" validate:"required"`
	Location string `json:"location" validate:"required"`
}

// convert to domain model.
func (r *CreateStockItemRequest) ToDomain() domain.StockItem {
	return domain.StockItem{
		UserID: domain.UserID(r.UserID),
		Sku: domain.SKU{
			ID: domain.SKUID(r.SkuID),			
		},
		Count:    r.Count,
		Price:    r.Price,
		Location: r.Location,
	}
}

type DeleteStockItemRequest struct {
	UserID int64  `json:"userID" validate:"required"`
	SkuID  uint32 `json:"skuID" validate:"required"`
}

type GetStockItemRequest struct {
	SkuID uint32 `json:"skuID" validate:"required"`
}

type GetStockItemResponse struct {
	SkuID    uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type FilterRequest struct {
	UserID      int64  `json:"userID" validate:"required"`
	Location    string `json:"location" validate:"required"`
	PageSize    int64  `json:"pageSize" validate:"required,gte=1"`
	CurrentPage int64  `json:"currentPage" validate:"required,gte=1"`
}

func (f *FilterRequest) ToDomain() domain.Filter {
	return domain.Filter{
		UserID:      domain.UserID(f.UserID),
		Location:    f.Location,
		PageSize:    f.PageSize,
		CurrentPage: f.CurrentPage,
	}
}

type StockItemResponse struct {
	SkuID    uint32 `json:"sku_id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type ListStockItemsResponse struct {
	Items      []StockItemResponse `json:"items"`
	TotalCount uint16              `json:"total_count"`
	PageNumber int64               `json:"page_number"`
}

func ToStockItemResponse(item domain.StockItem) StockItemResponse {
	return StockItemResponse{
		SkuID:    uint32(item.Sku.ID),
		Name:     item.Sku.Name,
		Type:     item.Sku.Type,
		Count:    item.Count,
		Price:    item.Price,
		Location: item.Location,
	}
}

func ToListResponse(p domain.PaginatedResponse[domain.StockItem]) ListStockItemsResponse {
	items := make([]StockItemResponse, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, ToStockItemResponse(item))
	}

	return ListStockItemsResponse{
		Items:      items,
		TotalCount: p.TotalCount,
		PageNumber: p.PageNumber,
	}
}
