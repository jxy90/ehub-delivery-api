package models

import (
	"context"
	"errors"

	"github.com/hublabs/ehub-delivery-api/factory"
)

var (
	DeliveryTableName     = "dlv_delivery"
	DeliveryItemTableName = "dlv_delivery_item"
)

func (d *Delivery) TableName() string {
	return DeliveryTableName
}

func (di *DeliveryItem) TableName() string {
	return DeliveryItemTableName
}

type Delivery struct {
	Id                 int64          `json:"id"`
	ShipmentLocationId int64          `json:"shipmentLocationId"`
	ReceiptLocationId  int64          `json:"receiptLocationId"`
	WaybillNo          string         `json:"waybillNo"`
	BoxNo              string         `json:"boxNo"`
	Status             string         `json:"status"`
	Items              []DeliveryItem `json:"items" xorm:"-"`
	Committed          `xorm:"extends"`
}

type DeliveryItem struct {
	Id          int64 `json:"id"`
	DeliveryId  int64 `json:"deliveryId"`
	SkuId       int64 `json:"skuId"`
	ShipmentQty int64 `json:"shipmentQty"`
	ReceiptQty  int64 `json:"receiptQty"`
	Committed   `xorm:"extends"`
}

func (d *Delivery) Create(ctx context.Context) error {
	if _, err := factory.
		DB(ctx).
		Table(DeliveryTableName).
		Insert(d); err != nil {
		return err
	}
	for i := range d.Items {
		d.Items[i].DeliveryId = d.Id
		d.Items[i].Committed = d.Committed
	}
	counts, err := DeliveryItem{}.insertDeliveryItems(ctx, d.Items)
	if err != nil {
		return err
	}
	if counts != int64(len(d.Items)) {
		return errors.New("fail to insert partial data")
	}
	return nil
}

func (DeliveryItem) insertDeliveryItems(ctx context.Context, items []DeliveryItem) (int64, error) {
	counts, err := factory.
		DB(ctx).
		Table(DeliveryItemTableName).
		Insert(items)
	if err != nil {
		return 0, err
	}
	return counts, nil
}
