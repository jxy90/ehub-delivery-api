package models

type StockCreateDto struct {
	LocationId int64                `json:"locationId"`
	Items      []StockItemCreateDto `json:"items"`
	CreatedBy  string               `json:"createdBy"`
}

type StockItemCreateDto struct {
	SkuId int64 `json:"skuId"`
	Qty   int64 `json:"qty"`
}
