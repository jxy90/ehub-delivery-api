package models

type StockCreateDto struct {
	LocationId int64                `json:"locationId" validate:"required"`
	Items      []StockItemCreateDto `json:"items" validate:"required,dive,required"`
	Type       string               `json:"type" validate:"required"` // todo 以后从token中获取
	CreatedBy  string               `json:"createdBy" validate:"required"`
}

type StockItemCreateDto struct {
	SkuId int64 `json:"skuId" validate:"required"`
	Qty   int64 `json:"qty" validate:"required"`
}
