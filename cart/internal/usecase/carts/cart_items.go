package carts

import (
	"cart/internal/domain"
	"cart/internal/usecase"
	"context"
	"errors"
)

type (
	// StockService interface represent stock service buisiness logic.
	StockService interface {
		GetStockItemBySKU(ctx context.Context, skuID domain.SkuID) (domain.StockItemBySKU, error)
	}
	// CartItemRepository interface represent cart items repository logic.
	CartItemRepository interface {
		SaveCartItem(ctx context.Context, cartItem domain.CartItem) error
		UpdateCartItem(ctx context.Context, cartItem domain.CartItem) error
		RemoveCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error
		RemoveAllCartItems(ctx context.Context, userID domain.UserID) error
		GetCartItemByUserID(ctx context.Context, userID domain.UserID, skuID domain.SkuID) (domain.CartItem, error)
	}
)

type cartServiceUseCase struct {
	StockService
	CartItemRepository
}

var _ usecase.CartItemUseCase = (*cartServiceUseCase)(nil)

func NewCartServiceUseCase(stockService StockService, cartItemRepo CartItemRepository) *cartServiceUseCase {
	return &cartServiceUseCase{
		StockService:       stockService,
		CartItemRepository: cartItemRepo,
	}
}

func (u *cartServiceUseCase) AddCartItem(ctx context.Context, cartItem domain.CartItem) error {
	stockItemBySKU, err := u.GetStockItemBySKU(ctx, cartItem.SkuID)
	if err != nil {
		return err
	}

	if cartItem.Count > stockItemBySKU.Count {
		return domain.ErrInSufficientStockCount
	}

	existingCartItem, err := u.GetCartItemByUserID(ctx, cartItem.UserID, cartItem.SkuID)
	if err != nil {
		if errors.Is(err, domain.ErrCartItemNotFound) {
			return u.SaveCartItem(ctx, cartItem)
		}

		return err
	}

	existingCartItem.Count = cartItem.Count

	return u.UpdateCartItem(ctx, existingCartItem)
}

func (u *cartServiceUseCase) DeleteCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error {
	return u.RemoveCartItem(ctx, userID, skuID)
}

func (u *cartServiceUseCase) ClearCartItems(ctx context.Context, userID domain.UserID) error {
	return u.RemoveAllCartItems(ctx, userID)
}
