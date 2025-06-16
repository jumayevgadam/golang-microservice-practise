package domain

import "errors"

// ErrSKUNotFound is used when sku not found by given skuID.
var ErrSKUNotFound = errors.New("sku not found")

// ErrStockItemNotFound is used when stock item not found.
var ErrStockItemNotFound = errors.New("stock item not found")
