package v1

import (
	"context"
	"errors"
	"stocks/internal/domain"
	"stocks/internal/usecase"
	pb "stocks/pkg/api/stocks"
	helper "stocks/pkg/httphelper"

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
	stockItemReq := fromGrpcStockItemReqToDomain(req)

	if err := helper.ValidateRequest(&stockItemReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "helper.ValidateRequest[CreateStockItemRequest]: %v", err)
	}

	err := s.stockUC.AddStockItem(ctx, stockItemReq)
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
	deleteStockItemReq := fromGrpcDeleteStockItemReqToDomain(req)

	if err := helper.ValidateRequest(&deleteStockItemReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "helper.ValidateRequest[DeleteStockItem]: %v", err)
	}

	err := s.stockUC.DeleteStockItem(ctx, deleteStockItemReq.UserID, deleteStockItemReq.Sku.ID)
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
	skuID := fromGrpcGetStockItemReqToDomain(req)

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
	filterReq := fromGrpcListStockItemsFilterToDomain(filter)

	if err := helper.ValidateRequest(&filterReq); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "helper.ValidateRequest[ListStockItemsByLocation]: %v", err)
	}

	listStockItems, err := s.stockUC.ListStockItems(ctx, filterReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return fromListStockItemsDomainToGrpc(listStockItems.Items, listStockItems.TotalCount, listStockItems.PageNumber), nil
}
