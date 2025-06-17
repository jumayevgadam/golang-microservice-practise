package domain

type CartItem struct {
	UserID UserID
	SkuID  SkuID
	Count  uint16
}

type ListCartItems struct {
	Items      []StockItemBySKU
	TotalPrice uint32
}
