package dto

type GetStockInfoInput struct {
	Sku uint32
}

type GetStockInfoOutput struct {
	Count int
}
