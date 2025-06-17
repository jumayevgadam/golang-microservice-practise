package usecase

import (
	"cart/internal/domain"
	"context"
)

type (
	CartItemUseCase interface {
		AddCartItem(ctx context.Context, cartItem domain.CartItem) error
		DeleteCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error
		ClearCartItems(ctx context.Context, userID domain.UserID) error
		ListCartItems(ctx context.Context, userID domain.UserID) (domain.ListCartItems, error)
	}
)
