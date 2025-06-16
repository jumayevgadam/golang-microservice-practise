package postgres

import (
	"cart/internal/domain"
	"cart/internal/usecase/carts"
	"cart/pkg/connection"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type cartServiceRepo struct {
	psqlDB connection.DB
}

var _ carts.CartItemRepository = (*cartServiceRepo)(nil)

func NewCartItemRepository(psqlDB connection.DB) *cartServiceRepo {
	return &cartServiceRepo{psqlDB: psqlDB}
}

func (c *cartServiceRepo) SaveCartItem(ctx context.Context, cartItem domain.CartItem) error {
	_, err := c.psqlDB.Exec(ctx, `
		INSERT INTO cart_items (user_id, sku, count)
		VALUES ($1, $2, $3)`,
		cartItem.UserID, cartItem.SkuID, cartItem.Count,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *cartServiceRepo) RemoveCartItem(ctx context.Context, userID domain.UserID, skuID domain.SkuID) error {
	_, err := c.psqlDB.Exec(ctx, `
		DELETE FROM cart_items
		WHERE user_id = $1 AND sku = $2`,
		userID, skuID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCartItemNotFound
		}

		return err
	}

	return nil
}

func (c *cartServiceRepo) UpdateCartItem(ctx context.Context, cartItem domain.CartItem) error {
	_, err := c.psqlDB.Exec(ctx, `
		UPDATE cart_items
		SET 
			count = COALESCE(NULLIF($1, 0), count),
			updated_at = NOW()
		WHERE user_id = $2 AND sku = $3`,
		cartItem.Count, cartItem.UserID, cartItem.SkuID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *cartServiceRepo) GetCartItemByUserID(ctx context.Context, userID domain.UserID, skuID domain.SkuID) (domain.CartItem, error) {
	var cartItemData CartItemData

	err := c.psqlDB.Get(ctx, &cartItemData, `
		SELECT user_id, sku, count, created_at, updated_at
		FROM cart_items
		WHERE user_id = $1 AND sku = $2`,
		userID, skuID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.CartItem{}, domain.ErrCartItemNotFound
		}

		return domain.CartItem{}, err
	}

	return cartItemData.ToDomain(), nil
}

func (c *cartServiceRepo) RemoveAllCartItems(ctx context.Context, userID domain.UserID) error {
	_, err := c.psqlDB.Exec(ctx, `
		DELETE FROM cart_items
		WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrCartItemNotFound
		}

		return err
	}

	return nil
}
