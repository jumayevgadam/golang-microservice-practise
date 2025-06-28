package carts

import (
	"cart/internal/domain"
	"cart/internal/usecase"
	"context"
)

type (
	// StockService interface represent stock service buisiness logic.
	StockService interface {
		GetStockItemBySKU(ctx context.Context, skuID domain.SkuID) (domain.StockItemBySKU, error)
	}
	// CartItemRepository interface represent cart items repository logic.
	CartItemRepository interface {
		SaveOrUpdateCartItem(ctx context.Context, cartItem domain.CartItem) error
		//UpdateCartItem(ctx context.Context, cartItem domain.CartItem) error
		RemoveCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error
		RemoveAllCartItems(ctx context.Context, userID domain.UserID) error
		GetCartItemByUserID(ctx context.Context, userID domain.UserID, skuID domain.SkuID) (domain.CartItem, error)
		ListCartItemsByUserID(ctx context.Context, userID domain.UserID) ([]domain.CartItem, error)
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

	return u.SaveOrUpdateCartItem(ctx, cartItem)
}

func (u *cartServiceUseCase) DeleteCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error {
	return u.RemoveCartItem(ctx, userID, skuID)
}

func (u *cartServiceUseCase) ClearCartItems(ctx context.Context, userID domain.UserID) error {
	return u.RemoveAllCartItems(ctx, userID)
}

func (u *cartServiceUseCase) ListCartItems(ctx context.Context, userID domain.UserID) (domain.ListCartItems, error) {
	var listCartItemsResponse domain.ListCartItems
	var totalPrice uint32

	listCartItems, err := u.ListCartItemsByUserID(ctx, userID)
	if err != nil {
		return domain.ListCartItems{}, err
	}
	// call service...
	stockItems := make([]domain.StockItemBySKU, 0, len(listCartItems))

	for _, listCartItem := range listCartItems {
		stockItem, err := u.GetStockItemBySKU(ctx, listCartItem.SkuID)
		if err != nil {
			continue
		}

		stockItems = append(stockItems, stockItem)
	}

	for _, stockItem := range stockItems {
		totalPrice += stockItem.Price * uint32(stockItem.Count)
	}

	listCartItemsResponse.Items = stockItems
	listCartItemsResponse.TotalPrice = totalPrice

	return listCartItemsResponse, nil
}
