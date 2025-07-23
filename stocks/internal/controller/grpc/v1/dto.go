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
