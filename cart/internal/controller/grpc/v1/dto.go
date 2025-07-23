package v1

import "cart/internal/domain"

type CreateCartItemRequest struct {
	UserID int64  `json:"userID" validate:"required"`
	SkuID  uint32 `json:"skuID" validate:"required"`
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
	UserID int64  `json:"userID" validate:"required"`
	SkuID  uint32 `json:"skuID" validate:"required"`
}

type ClearCartItemRequest struct {
	UserID int64 `json:"userID" validate:"required"`
}

type ListCartItemsRequest struct {
	UserID int64 `json:"userID" validate:"required"`
}
