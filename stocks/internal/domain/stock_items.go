package domain

// StockItem represent stock's items domain.
type StockItem struct {
	UserID   UserID
	Sku      SKU
	Count    uint16
	Price    uint32
	Location string
}
