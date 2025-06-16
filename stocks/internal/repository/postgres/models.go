package postgres

import (
	"stocks/internal/domain"
	"time"
)

type SKU struct {
	SkuID uint32 `db:"sku_id"`
	Name  string `db:"name"`
	Type  string `db:"type"`
}

func (s *SKU) ToDomain() domain.SKU {
	return domain.SKU{
		ID:   domain.SKUID(s.SkuID),
		Name: s.Name,
		Type: s.Type,
	}
}

type StockItemData struct {
	UserID    int64     `db:"user_id"`
	SkuID     uint32    `db:"sku"`
	Count     uint16    `db:"count"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Price     uint32    `db:"price"`
	Location  string    `db:"location"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s *StockItemData) ToDomain() domain.StockItem {
	return domain.StockItem{
		UserID: domain.UserID(s.UserID),
		Sku: domain.SKU{
			ID:   domain.SKUID(s.SkuID),
			Name: s.Name,
			Type: s.Type,
		},
		Count:    s.Count,
		Price:    s.Price,
		Location: s.Location,
	}
}
