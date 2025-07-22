package v1

import (
	"stocks/internal/domain"
	"stocks/pkg/api/stocks"
)

func stockItemReqToDomain(req *stocks.CreateStockItemRequest) domain.StockItem {
	return domain.StockItem{
		UserID: domain.UserID(req.UserId),
		Sku: domain.SKU{
			ID: domain.SKUID(req.SkuId),
		},
		Count:    uint16(req.Count),
		Price:    req.Price,
		Location: req.Location,
	}
}
