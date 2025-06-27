package postgres

import (
	"context"
	"errors"
	"stocks/internal/domain"
	"stocks/internal/usecase/stocks"
	"stocks/pkg/connection"

	"github.com/jackc/pgx/v5"
)

var _ stocks.StockServiceRepository = (*stockServiceRepository)(nil)

type stockServiceRepository struct {
	psqlDB connection.DB
}

func NewStockServiceRepository(psqlDB connection.DB) *stockServiceRepository {
	return &stockServiceRepository{psqlDB: psqlDB}
}

func (s *stockServiceRepository) SaveOrUpdateStockItem(ctx context.Context, stockItem domain.StockItem) error {
	_, err := s.psqlDB.Exec(ctx, `
		INSERT INTO stock_items (user_id, sku_id, count, price, location)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, sku_id) DO UPDATE SET
			count = stock_items.count + EXCLUDED.count,
			price = COALESCE(NULLIF(EXCLUDED.price, 0), stock_items.price),
			location = COALESCE(NULLIF(EXCLUDED.location, ''), stock_items.location),
			updated_at = NOW()`,
		stockItem.UserID, stockItem.Sku.ID, stockItem.Count,
		stockItem.Price, stockItem.Location,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *stockServiceRepository) GetStockItem(ctx context.Context, userID domain.UserID, skuID domain.SKUID) (domain.StockItem, error) {
	var stockItemData StockItemData

	err := s.psqlDB.Get(ctx, &stockItemData, `
		SELECT si.user_id, s.sku_id, si.count, s.name, s.type, si.price, si.location, si.created_at, si.updated_at
		FROM stock_items si 
		LEFT JOIN sku s ON s.sku_id = si.sku_id
		WHERE si.user_id = $1 AND si.sku_id = $2`,
		userID,
		skuID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.StockItem{}, domain.ErrStockItemNotFound
		}

		return domain.StockItem{}, err
	}

	return stockItemData.ToDomain(), nil
}

func (s *stockServiceRepository) UpdateStockItem(ctx context.Context, stockItem domain.StockItem) error {
	_, err := s.psqlDB.Exec(ctx, `
		UPDATE stock_items
		SET	
			count = COALESCE(NULLIF($1, 0), count),
			price = COALESCE(NULLIF($2, 0), price),
			location = COALESCE(NULLIF($3, ''), location),
			updated_at = NOW()
		WHERE user_id = $4 AND sku_id = $5`,
		stockItem.Count, stockItem.Price, stockItem.Location,
		stockItem.UserID, stockItem.Sku.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *stockServiceRepository) DeleteStockItemFromStorage(ctx context.Context, userID domain.UserID, skuID domain.SKUID) error {
	_, err := s.psqlDB.Exec(ctx, `
		DELETE FROM stock_items
		WHERE user_id = $1 AND sku_id = $2`,
		userID, skuID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrStockItemNotFound
		}

		return err
	}

	return nil
}

func (s *stockServiceRepository) GetStockItemBySku(ctx context.Context, skuID domain.SKUID) (domain.StockItem, error) {
	var stockItemData StockItemData

	err := s.psqlDB.Get(ctx, &stockItemData, `
		SELECT si.user_id, s.sku_id, s.name, s.type, si.count, si.price, si.location, si.created_at, si.updated_at
		FROM stock_items si
		INNER JOIN sku s ON s.sku_id = si.sku_id
		WHERE si.sku_id = $1`,
		skuID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.StockItem{}, domain.ErrStockItemNotFound
		}

		return domain.StockItem{}, err
	}

	return stockItemData.ToDomain(), nil
}

func (s *stockServiceRepository) CountStockItems(ctx context.Context, userID domain.UserID, location string) (uint16, error) {
	var stockItemsCount uint16

	err := s.psqlDB.Get(ctx, &stockItemsCount, `
		SELECT COUNT(user_id) 
		FROM stock_items
		WHERE user_id = $1 AND location = $2`,
		userID, location,
	)
	if err != nil {
		return 0, err
	}

	return stockItemsCount, nil
}

func (s *stockServiceRepository) ListStockItemsByLocation(ctx context.Context, filter domain.Filter) ([]domain.StockItem, error) {
	var stockItemsData []StockItemData

	offset := (filter.CurrentPage - 1) * filter.PageSize

	err := s.psqlDB.Select(ctx, &stockItemsData, `
		SELECT si.user_id, s.sku_id, s.name, s.type, si.count, si.price, si.location, si.created_at, si.updated_at
		FROM stock_items si
		LEFT JOIN sku s ON s.sku_id = si.sku_id
		WHERE user_id = $1 AND location = $2
		OFFSET $3 LIMIT $4`,
		filter.UserID,
		filter.Location,
		offset,
		filter.PageSize,
	)
	if err != nil {
		return nil, err
	}

	stockItems := make([]domain.StockItem, 0, len(stockItemsData))
	for _, stockItem := range stockItemsData {
		stockItems = append(stockItems, stockItem.ToDomain())
	}

	return stockItems, nil
}
