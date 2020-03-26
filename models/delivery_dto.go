package models

import (
	"errors"
)

type ShipmentDto struct {
	ShipmentLocationId int64             `json:"shipmentLocationId"`
	ReceiptLocationId  int64             `json:"receiptLocationId"`
	WaybillNo          string            `json:"waybillNo"`
	BoxNo              string            `json:"boxNo"`
	PlatformOrderId    string            `json:"platformOrderId"`
	Items              []ShipmentItemDto `json:"items"`
	Type               string            `json:"type"` // todo 以后从token中获取
	CreatedBy          string            `json:"createdBy"`
}

type ShipmentItemDto struct {
	SkuId int64 `json:"skuId"`
	Qty   int64 `json:"qty"`
}

func (d ShipmentDto) translateToDelivery() (DeliveryProcessor, error) {
	switch d.Type {
	case Store.Code:
		var items []DeliveryItemForStore
		for _, i := range d.Items {
			item := DeliveryItemForStore{
				SkuId:       i.SkuId,
				ShipmentQty: i.Qty,
				ReceiptQty:  0,
				Committed:   Committed{}.newCommitted(d.CreatedBy),
			}
			items = append(items, item)
		}
		return &DeliveryForStore{
			ShipmentLocationId: d.ShipmentLocationId,
			ReceiptLocationId:  d.ReceiptLocationId,
			WaybillNo:          d.WaybillNo,
			BoxNo:              d.BoxNo,
			Status:             Shipment.Code,
			Items:              items,
			Committed:          Committed{}.newCommitted(d.CreatedBy),
		}, nil
	case Plant.Code:
		var items []DeliveryItemForPlant
		for _, i := range d.Items {
			item := DeliveryItemForPlant{
				SkuId:       i.SkuId,
				ShipmentQty: i.Qty,
				ReceiptQty:  0,
				Committed:   Committed{}.newCommitted(d.CreatedBy),
			}
			items = append(items, item)
		}
		return &DeliveryForPlant{
			ShipmentLocationId: d.ShipmentLocationId,
			ReceiptLocationId:  d.ReceiptLocationId,
			WaybillNo:          d.WaybillNo,
			BoxNo:              d.BoxNo,
			PlatformOrderId:    d.PlatformOrderId,
			Status:             Shipment.Code,
			Items:              items,
			Committed:          Committed{}.newCommitted(d.CreatedBy),
		}, nil
	}
	return nil, errors.New("not support type")
}

type ReceiptDto struct {
	DeliveryId int64            `json:"deliveryId"`
	Type       string           `json:"type"` // todo 以后从token中获取
	Items      []ReceiptItemDto `json:"items"`
	UpdatedBy  string           `json:"updatedBy"`
}

type ReceiptItemDto struct {
	SkuId int64 `json:"skuId"`
	Qty   int64 `json:"qty"`
}

func (d ReceiptDto) translateToDelivery() (DeliveryProcessor, error) {
	switch d.Type {
	case Store.Code:
		var items []DeliveryItemForStore
		for _, item := range d.Items {
			di := DeliveryItemForStore{
				DeliveryId: d.DeliveryId,
				SkuId:      item.SkuId,
				ReceiptQty: item.Qty,
				Committed: Committed{
					UpdatedBy: d.UpdatedBy,
				},
			}
			items = append(items, di)
		}
		return &DeliveryForStore{
			Id:     d.DeliveryId,
			Status: Receipt.Code,
			Items:  items,
			Committed: Committed{
				UpdatedBy: d.UpdatedBy,
			},
		}, nil
	case Plant.Code:
		var items []DeliveryItemForPlant
		for _, item := range d.Items {
			di := DeliveryItemForPlant{
				DeliveryId: d.DeliveryId,
				SkuId:      item.SkuId,
				ReceiptQty: item.Qty,
				Committed: Committed{
					UpdatedBy: d.UpdatedBy,
				},
			}
			items = append(items, di)
		}
		return &DeliveryForPlant{
			Id:     d.DeliveryId,
			Status: Receipt.Code,
			Items:  items,
			Committed: Committed{
				UpdatedBy: d.UpdatedBy,
			},
		}, nil
	}
	return nil, errors.New("not support type")
}
