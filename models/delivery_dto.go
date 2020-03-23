package models

type DeliveryCreateDto struct {
	ShipmentLocationId int64                   `json:"shipmentLocationId"`
	ReceiptLocationId  int64                   `json:"receiptLocationId"`
	WaybillNo          string                  `json:"waybillNo"`
	BoxNo              string                  `json:"boxNo"`
	Items              []DeliveryItemCreateDto `json:"items"`
	CreatedBy          string                  `json:"createdBy"`
}

type DeliveryItemCreateDto struct {
	SkuId int64 `json:"skuId"`
	Qty   int64 `json:"qty"`
}

type DeliveryReceiptDto struct {
	DeliveryId int64  `json:"deliveryId"`
	UpdatedBy  string `json:"updatedBy"`
}
