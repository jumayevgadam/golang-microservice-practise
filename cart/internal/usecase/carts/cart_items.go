package carts

import (
	"cart/internal/domain"
	"cart/internal/kafka"
	"cart/internal/usecase"
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := otel.Tracer(os.Getenv("SERVICE_NAME")).Start(ctx, "CartServiceUseCase.AddCartItem")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", fmt.Sprintf("%d", cartItem.UserID)),
		attribute.String("sku_id", fmt.Sprintf("%d", cartItem.SkuID)),
		attribute.Int64("count", int64(cartItem.Count)),
	)

	stockItemBySKU, err := u.GetStockItemBySKU(ctx, cartItem.SkuID)
	if err != nil {
		span.SetAttributes(attribute.String("error.message", err.Error()))
		return err
	}

	// prepare cart item addedpayload for producing event.
	payload := kafka.CartItemAddedPayload{
		CartID: fmt.Sprintf("%d", cartItem.UserID), // assuming userID is cartID.
		SKU:    uint32(cartItem.SkuID),
		Count:  cartItem.Count,
		Status: "success",
	}

	if cartItem.Count > stockItemBySKU.Count {
		u.KafkaProducer.ProduceCartItemFailed(ctx, kafka.CartItemFailedPayload{
			CartID: fmt.Sprintf("%d", cartItem.UserID),
			SKU:    uint32(cartItem.SkuID),
			Count:  cartItem.Count,
			Status: "failed",
			Reason: "not enough stock",
		})

		span.SetAttributes(attribute.String("error.message", domain.ErrInSufficientStockCount.Error()))

		return domain.ErrInSufficientStockCount
	}

	u.KafkaProducer.ProduceCartItemAdded(ctx, payload)

	return u.SaveOrUpdateCartItem(ctx, cartItem)
}

func (u *cartServiceUseCase) DeleteCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error {
	ctx, span := otel.Tracer(os.Getenv("SERVICE_NAME")).Start(ctx, "CartServiceUseCase.DeleteCartItem")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", fmt.Sprintf("%d", userID)),
		attribute.String("sku_id", fmt.Sprintf("%d", skuID)),
	)

	err := u.RemoveCartItem(ctx, userID, skuID)
	if err != nil {
		span.SetAttributes(attribute.String("error.message", err.Error()))
		return err
	}

	return nil
}

func (u *cartServiceUseCase) ClearCartItems(ctx context.Context, userID domain.UserID) error {
	ctx, span := otel.Tracer(os.Getenv("SERVICE_NAME")).Start(ctx, "CartServiceUseCase.ClearCartItems")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", fmt.Sprintf("%d", userID)),
	)

	err := u.RemoveAllCartItems(ctx, userID)
	if err != nil {
		span.SetAttributes(attribute.String("error.message", err.Error()))
		return err
	}

	return nil
}

func (u *cartServiceUseCase) ListCartItems(ctx context.Context, userID domain.UserID) (domain.ListCartItems, error) {
	ctx, span := otel.Tracer(os.Getenv("SERVICE_NAME")).Start(ctx, "CartServiceUseCase.ListCartItems")
	defer span.End()

	span.SetAttributes(
		attribute.String("user_id", fmt.Sprintf("%d", userID)),
	)

	var listCartItemsResponse domain.ListCartItems
	var totalPrice uint32

	listCartItems, err := u.ListCartItemsByUserID(ctx, userID)
	if err != nil {
		span.SetAttributes(attribute.String("error.message", err.Error()))
		return domain.ListCartItems{}, err
	}
	// call service...
	stockItems := make([]domain.StockItemBySKU, 0, len(listCartItems))

	for _, listCartItem := range listCartItems {
		stockItem, err := u.GetStockItemBySKU(ctx, listCartItem.SkuID)
		if err != nil {
			span.SetAttributes(attribute.String("error.message", err.Error()))
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
