package domain

import "errors"

// ErrInSufficientStockCount is returned when we got unsufficient stock count from cart service.
var ErrInSufficientStockCount = errors.New("insufficient stock count")

// ErrCartItemNotFound is returned when no rows in result set for cartItem.
var ErrCartItemNotFound = errors.New("cart item not found")
