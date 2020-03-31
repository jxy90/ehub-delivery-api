package models

import (
	"errors"
)

type ShipmentDto struct {
	ShipmentLocationId int64             `json:"shipmentLocationId" validate:"required"`
	ReceiptLocationId  int64             `json:"receiptLocationId" validate:"required"`
	WaybillNo          string            `json:"waybillNo" validate:"required"`
	BoxNo              string            `json:"boxNo" validate:"required"`
	PlatformOrderId    string            `json:"platformOrderId"`
	Items              []ShipmentItemDto `json:"items" validate:"required,dive,required"`
	Type               string            `json:"type" validate:"required"`      // todo 以后从token中获取
	ShippedBy          string            `json:"shippedBy" validate:"required"` // todo 以后从token中获取
}

type ShipmentItemDto struct {
	SkuId int64 `json:"skuId" validate:"required"`
	Qty   int64 `json:"qty" validate:"required"`
}

func (d ShipmentDto) translateToDeliveryByType() (DeliveryProcessor, error) {
	switch d.Type {
	case Store.Code:
		var items []DeliveryItemForStore
		for _, i := range d.Items {
			item := DeliveryItemForStore{
				SkuId:       i.SkuId,
				ShipmentQty: i.Qty,
				ReceiptQty:  0,
				Committed:   Committed{}.newCommitted(d.ShippedBy),
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
			Committed:          Committed{}.newCommitted(d.ShippedBy),
		}, nil
	case Plant.Code:
		var items []DeliveryItemForPlant
		for _, i := range d.Items {
			item := DeliveryItemForPlant{
				SkuId:       i.SkuId,
				ShipmentQty: i.Qty,
				ReceiptQty:  0,
				Committed:   Committed{}.newCommitted(d.ShippedBy),
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
			Committed:          Committed{}.newCommitted(d.ShippedBy),
		}, nil
	}
	return nil, errors.New("not support type, supported types are [store, plant]")
}

type ReceiptDto struct {
	DeliveryId        int64            `json:"deliveryId" validate:"required"`
	ReceiptLocationId int64            `json:"receiptLocationId" validate:"required"`
	Type              string           `json:"type" validate:"required"` // todo 以后从token中获取
	Items             []ReceiptItemDto `json:"items" validate:"required,dive,required"`
	ReceiptedBy       string           `json:"receiptedBy" validate:"required"` // todo 以后从token中获取
}

type ReceiptItemDto struct {
	SkuId int64 `json:"skuId" validate:"required"`
	Qty   int64 `json:"qty" validate:"required"`
}

func (d ReceiptDto) translateToDeliveryByType() (DeliveryProcessor, error) {
	switch d.Type {
	case Store.Code:
		var items []DeliveryItemForStore
		for _, item := range d.Items {
			di := DeliveryItemForStore{
				DeliveryId: d.DeliveryId,
				SkuId:      item.SkuId,
				ReceiptQty: item.Qty,
				Committed: Committed{
					UpdatedBy: d.ReceiptedBy,
				},
			}
			items = append(items, di)
		}
		return &DeliveryForStore{
			Id:                d.DeliveryId,
			ReceiptLocationId: d.ReceiptLocationId,
			Status:            Receipt.Code,
			Items:             items,
			Committed: Committed{
				UpdatedBy: d.ReceiptedBy,
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
					UpdatedBy: d.ReceiptedBy,
				},
			}
			items = append(items, di)
		}
		return &DeliveryForPlant{
			Id:                d.DeliveryId,
			ReceiptLocationId: d.ReceiptLocationId,
			Status:            Receipt.Code,
			Items:             items,
			Committed: Committed{
				UpdatedBy: d.ReceiptedBy,
			},
		}, nil
	}
	return nil, errors.New("not support type, supported types are [store, plant]")
}
