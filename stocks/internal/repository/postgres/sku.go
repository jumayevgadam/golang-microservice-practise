package postgres

import (
	"context"
	"errors"
	"stocks/internal/domain"
	"stocks/internal/usecase/stocks"
	"stocks/pkg/connection"

	"github.com/jackc/pgx/v5"
)

var _ stocks.SKURepository = (*skuRepository)(nil)

type skuRepository struct {
	psqlDB connection.DB
}

func NewSKURepository(psqlDB connection.DB) *skuRepository {
	return &skuRepository{psqlDB: psqlDB}
}

func (s *skuRepository) GetSKUByID(ctx context.Context, skuID domain.SKUID) (domain.SKU, error) {
	var sku SKU

	err := s.psqlDB.Get(ctx, &sku, `
		SELECT sku_id, name, type FROM sku WHERE sku_id = $1`,
		skuID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.SKU{}, domain.ErrSKUNotFound
		}

		return domain.SKU{}, err
	}

	return sku.ToDomain(), nil
}
