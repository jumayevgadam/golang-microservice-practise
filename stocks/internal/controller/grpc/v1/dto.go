package v1

import (
	"stocks/internal/domain"
	"stocks/pkg/api/stocks"
)

func fromGrpcStockItemReqToDomain(req *stocks.CreateStockItemRequest) domain.StockItem {
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

func fromGrpcDeleteStockItemReqToDomain(req *stocks.DeleteStockItemRequest) domain.StockItem {
	return domain.StockItem{
		UserID: domain.UserID(req.UserId),
		Sku: domain.SKU{
			ID: domain.SKUID(req.SkuId),
		},
	}
}

func fromGrpcGetStockItemReqToDomain(req *stocks.GetStockItemRequest) domain.SKUID {
	return domain.SKUID(req.SkuId)
}

func fromStockItemDomainToGrpc(stockItem domain.StockItem) *stocks.StockItemResponse {
	return &stocks.StockItemResponse{
		SkuId:    uint32(stockItem.Sku.ID),
		Name:     stockItem.Sku.Name,
		Type:     stockItem.Sku.Type,
		Count:    uint32(stockItem.Count),
		Price:    stockItem.Price,
		Location: stockItem.Location,
	}
}

func fromGrpcListStockItemsFilterToDomain(filter *stocks.FilterRequest) domain.Filter {
	return domain.Filter{
		UserID:      domain.UserID(filter.UserId),
		Location:    filter.Location,
		PageSize:    filter.PageSize,
		CurrentPage: filter.CurrentPage,
	}
}

func fromListStockItemsDomainToGrpc(stockItems []domain.StockItem, totalCount uint16, pageNumber int64) *stocks.ListStockItemsResponse {
	stockItemResponses := make([]*stocks.StockItemResponse, 0, len(stockItems))

	for _, stockItem := range stockItems {
		stockItemResponses = append(stockItemResponses, fromStockItemDomainToGrpc(stockItem))
	}

	return &stocks.ListStockItemsResponse{
		Items:      stockItemResponses,
		TotalCount: uint32(totalCount),
		PageNumber: pageNumber,
	}
}
