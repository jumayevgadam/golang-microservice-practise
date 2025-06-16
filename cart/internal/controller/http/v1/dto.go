package v1

import "cart/internal/domain"

type CreateCartItemRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	SkuID  uint32 `json:"sku_id" validate:"required"`
	Count  uint16 `json:"count" validate:"required"`
}

func (c *CreateCartItemRequest) ToDomain() domain.CartItem {
	return domain.CartItem{
		UserID: domain.UserID(c.UserID),
		SkuID:  domain.SkuID(c.SkuID),
		Count:  c.Count,
	}
}

type DeleteCartItemRequest struct {
	UserID int64  `json:"user_id" validate:"required"`
	SkuID  uint32 `json:"sku_id" validate:"required"`
}

type ClearCartItemRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type ListCartItemsRequest struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type CartItemResponse struct {
	SkuID uint32 `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint16 `json:"count"`
}

type ListCartItemsResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}

func ToCartItemResponse(item domain.StockItemBySKU) CartItemResponse {
	return CartItemResponse{
		SkuID: uint32(item.SKuID),
		Name:  item.Name,
		Price: item.Price,
		Count: item.Count,
	}
}

func ToListResponse(p domain.ListCartItems) ListCartItemsResponse {
	items := make([]CartItemResponse, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, ToCartItemResponse(item))
	}

	return ListCartItemsResponse{
		Items:      items,
		TotalPrice: p.TotalPrice,
	}
}
