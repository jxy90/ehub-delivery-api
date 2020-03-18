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

func (dto DeliveryCreateDto) Translate() (Delivery, error) {
	var items []DeliveryItem
	for _, i := range dto.Items {
		item := DeliveryItem{
			SkuId:       i.SkuId,
			ShipmentQty: i.Qty,
			ReceiptQty:  0,
			Committed:   Committed{}.newCommitted(dto.CreatedBy),
		}
		items = append(items, item)
	}
	return Delivery{
		ShipmentLocationId: dto.ShipmentLocationId,
		ReceiptLocationId:  dto.ReceiptLocationId,
		WaybillNo:          dto.WaybillNo,
		BoxNo:              dto.BoxNo,
		Status:             Shipment.Code,
		Items:              items,
		Committed:          Committed{}.newCommitted(dto.CreatedBy),
	}, nil
}
