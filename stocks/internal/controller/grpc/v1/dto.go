package v1

import (
	dtoModels "stocks/internal/controller/http/v1"
	"stocks/internal/domain"
	"stocks/pkg/api/stocks"
	helper "stocks/pkg/httphelper"
)

func fromGrpcStockItemReqToDomain(req *stocks.CreateStockItemRequest) (domain.StockItem, error) {
	createStockItemReq := dtoModels.CreateStockItemRequest{
		SkuID:    req.SkuId,
		UserID:   req.UserId,
		Count:    uint16(req.Count),
		Price:    req.Price,
		Location: req.Location,
	}

	if err := helper.ValidateRequest(&createStockItemReq); err != nil {
		return domain.StockItem{}, err
	}

	return createStockItemReq.ToDomain(), nil
}

func fromGrpcDeleteStockItemReqToDomain(req *stocks.DeleteStockItemRequest) (domain.StockItem, error) {
	deleteStockItemReq := dtoModels.DeleteStockItemRequest{
		UserID: req.UserId,
		SkuID:  req.SkuId,
	}

	if err := helper.ValidateRequest(&deleteStockItemReq); err != nil {
		return domain.StockItem{}, err
	}

	return domain.StockItem{
		UserID: domain.UserID(deleteStockItemReq.UserID),
		Sku: domain.SKU{
			ID: domain.SKUID(deleteStockItemReq.SkuID),
		},
	}, nil
}

func fromGrpcGetStockItemReqToDomain(req *stocks.GetStockItemRequest) (domain.SKUID, error) {
	getStockItemReq := dtoModels.GetStockItemRequest{
		SkuID: req.SkuId,
	}

	if err := helper.ValidateRequest(&getStockItemReq); err != nil {
		return domain.SKUID(0), err
	}

	return domain.SKUID(getStockItemReq.SkuID), nil
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

func fromGrpcListStockItemsFilterToDomain(filter *stocks.FilterRequest) (domain.Filter, error) {
	filterRequest := dtoModels.FilterRequest{
		UserID:      filter.UserId,
		Location:    filter.Location,
		PageSize:    filter.PageSize,
		CurrentPage: filter.CurrentPage,
	}

	if err := helper.ValidateRequest(&filterRequest); err != nil {
		return domain.Filter{}, err
	}

	return filterRequest.ToDomain(), nil
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
