package v1

import (
	"context"
	"errors"
	"stocks/internal/domain"
	"stocks/internal/usecase"
	pb "stocks/pkg/api/stocks"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	stockItemNotFound = "stock item not found"
)

type StockGRPCHandler struct {
	pb.UnimplementedStocksServiceServer
	stockUC usecase.StockServiceUseCase
}

func NewStockGRPCHandler(stockUC usecase.StockServiceUseCase) *StockGRPCHandler {
	return &StockGRPCHandler{stockUC: stockUC}
}

func (s *StockGRPCHandler) AddStockItem(ctx context.Context, req *pb.CreateStockItemRequest) (*pb.GeneralResponse, error) {
	stockItemReq, err := fromGrpcStockItemReqToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.stockUC.AddStockItem(ctx, stockItemReq)
	if err != nil {
		if errors.Is(err, domain.ErrSKUNotFound) {
			return nil, status.Error(codes.NotFound, "SKU not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Success: true,
		Message: "Stock item added successfully",
	}, nil
}

func (s *StockGRPCHandler) DeleteStockItem(ctx context.Context, req *pb.DeleteStockItemRequest) (*pb.GeneralResponse, error) {
	deleteStockItemReq, err := fromGrpcDeleteStockItemReqToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.stockUC.DeleteStockItem(ctx, deleteStockItemReq.UserID, deleteStockItemReq.Sku.ID)
	if err != nil {
		if errors.Is(err, domain.ErrStockItemNotFound) {
			return nil, status.Error(codes.NotFound, stockItemNotFound)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Success: true,
		Message: "stock item successfully removed",
	}, nil
}

func (s *StockGRPCHandler) GetStockItemBySKU(ctx context.Context, req *pb.GetStockItemRequest) (*pb.StockItemResponse, error) {
	skuID, err := fromGrpcGetStockItemReqToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	stockItem, err := s.stockUC.GetStockItemBySKU(ctx, skuID)
	if err != nil {
		if errors.Is(err, domain.ErrStockItemNotFound) {
			return nil, status.Error(codes.NotFound, stockItemNotFound)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return fromStockItemDomainToGrpc(stockItem), nil
}

func (s *StockGRPCHandler) ListStockItemsByLocation(ctx context.Context, filter *pb.FilterRequest) (*pb.ListStockItemsResponse, error) {
	filterReq, err := fromGrpcListStockItemsFilterToDomain(filter)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	listStockItems, err := s.stockUC.ListStockItems(ctx, filterReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return fromListStockItemsDomainToGrpc(listStockItems.Items, listStockItems.TotalCount, listStockItems.PageNumber), nil
}
