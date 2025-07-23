package v1

import (
	"cart/internal/domain"
	"cart/internal/usecase"
	pb "cart/pkg/api/cart"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cartItemNotFound = "cart item not found"
)

type CartGRPCHandler struct {
	pb.UnimplementedCartServiceServer
	cartUC usecase.CartItemUseCase
}

func NewCartGRPCHandler(cartUC usecase.CartItemUseCase) *CartGRPCHandler {
	return &CartGRPCHandler{cartUC: cartUC}
}

func (c *CartGRPCHandler) AddCartItem(ctx context.Context, req *pb.CreateCartItemRequest) (*pb.GeneralResponse, error) {
	cartItemReq, err := fromGrpcCreateCartItemReqToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.cartUC.AddCartItem(ctx, cartItemReq)
	if err != nil {
		if errors.Is(err, domain.ErrInSufficientStockCount) {
			return nil, status.Error(codes.InvalidArgument, "insufficient stock count")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Success: true,
		Message: "cart item added successfully",
	}, nil
}

func (c *CartGRPCHandler) DeleteCartItem(ctx context.Context, req *pb.RemoveCartItemRequest) (*pb.GeneralResponse, error) {
	deleteCartItemReq, err := fromGrpcDeleteCartItemReqToDomain(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = c.cartUC.DeleteCartItem(ctx, deleteCartItemReq.UserID, deleteCartItemReq.SkuID)
	if err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			return nil, status.Error(codes.NotFound, cartItemNotFound)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Success: true,
		Message: "cart item deleted successfully",
	}, nil
}

func (c *CartGRPCHandler) ClearCartItems(ctx context.Context, req *pb.ClearCartItemRequest) (*pb.GeneralResponse, error) {
	userID := domain.UserID(req.UserId)
	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be zero")
	}

	err := c.cartUC.ClearCartItems(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			return nil, status.Error(codes.NotFound, "cart item not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GeneralResponse{
		Success: true,
		Message: "cart items cleared succcessfully",
	}, nil
}

func (c *CartGRPCHandler) ListCartItems(ctx context.Context, req *pb.ListCartItemsRequest) (*pb.ListCartItemsResponse, error) {
	userID := domain.UserID(req.UserId)
	if userID == 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID cannot be zero")
	}

	listCartItems, err := c.cartUC.ListCartItems(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return fromListStockItemsDomainToGrpc(listCartItems), nil
}
