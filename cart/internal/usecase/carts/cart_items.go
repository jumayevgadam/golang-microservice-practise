package carts

import (
	"cart/internal/domain"
	"cart/internal/kafka"
	"cart/internal/usecase"
	"context"
	"fmt"
)

//go:generate mkdir -p mock
//go:generate minimock -o ./mock/ -s .go -g
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
	KafkaProducer kafka.CartEventProducer
}

var _ usecase.CartItemUseCase = (*cartServiceUseCase)(nil)

func NewCartServiceUseCase(
	stockService StockService,
	cartItemRepo CartItemRepository,
	kafkaProducer kafka.CartEventProducer,
) *cartServiceUseCase {
	return &cartServiceUseCase{
		StockService:       stockService,
		CartItemRepository: cartItemRepo,
		KafkaProducer:      kafkaProducer,
	}
}

func (u *cartServiceUseCase) AddCartItem(ctx context.Context, cartItem domain.CartItem) error {
	stockItemBySKU, err := u.GetStockItemBySKU(ctx, cartItem.SkuID)
	if err != nil {
		return err
	}

	// prepare cart item addedpayload for producing event.
	payload := kafka.CartItemAddedPayload{
		CartID: fmt.Sprintf("%d", cartItem.UserID), // assuming userID is cartID.
		SKU:    uint32(cartItem.SkuID),
		Count:  uint16(cartItem.Count),
		Status: "success",
	}

	if cartItem.Count > stockItemBySKU.Count {
		u.KafkaProducer.ProduceCartItemFailed(ctx, kafka.CartItemFailedPayload{
			CartID: fmt.Sprintf("%d", cartItem.UserID),
			SKU:    uint32(cartItem.SkuID),
			Count:  uint16(cartItem.Count),
			Status: "failed",
			Reason: "not enough stock",
		})
		return domain.ErrInSufficientStockCount
	}

	u.KafkaProducer.ProduceCartItemAdded(ctx, payload)

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
