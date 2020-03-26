package models

type StockCreateDto struct {
	LocationId int64                `json:"locationId"`
	Items      []StockItemCreateDto `json:"items"`
	Type       string               `json:"type"` // todo 以后从token中获取
	CreatedBy  string               `json:"createdBy"`
}

type StockItemCreateDto struct {
	SkuId int64 `json:"skuId"`
	Qty   int64 `json:"qty"`
}
