package usecase

import (
	"context"
	"stocks/internal/domain"
)

type (
	// StockServiceUseCase represent stock service usecase methods.
	StockServiceUseCase interface {
		AddStockItem(ctx context.Context, stockItem domain.StockItem) error
		DeleteStockItem(ctx context.Context, userID domain.UserID, skuID domain.SKUID) error
		GetStockItemBySKU(ctx context.Context, skuID domain.SKUID) (domain.StockItem, error)
		ListStockItems(ctx context.Context, filter domain.Filter) (domain.PaginatedResponse[domain.StockItem], error)
	}
)
