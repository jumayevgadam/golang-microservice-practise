package domain

type StockItemBySKU struct {
	SKuID SkuID
	Name  string
	Price uint32
	Count uint16
}
