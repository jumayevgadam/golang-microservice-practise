package postgres

import (
	"cart/internal/domain"
	"time"
)

type CartItemData struct {
	UserID    int64     `db:"user_id"`
	SkuID     uint32    `db:"sku"`
	Count     uint16    `db:"count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *CartItemData) ToDomain() domain.CartItem {
	return domain.CartItem{
		UserID: domain.UserID(c.UserID),
		SkuID:  domain.SkuID(c.SkuID),
		Count:  c.Count,
	}
}
