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
	ShipmentLocationId int64          `json:"shipmentLocationId" xorm:"index unique(delivery)"`
	ReceiptLocationId  int64          `json:"receiptLocationId" xorm:"index unique(delivery)"`
	WaybillNo          string         `json:"waybillNo" xorm:"index unique(delivery)"`
	BoxNo              string         `json:"boxNo" xorm:"index unique(delivery)"`
	Status             string         `json:"status"`
	Items              []DeliveryItem `json:"items" xorm:"-"`
	PlatformOrderId    string         `json:"platformOrderId"`
	Committed          `xorm:"extends"`
}

type DeliveryItem struct {
	Id          int64 `json:"id"`
	DeliveryId  int64 `json:"deliveryId" xorm:"index unique(delivery_item)"`
	SkuId       int64 `json:"skuId" xorm:"index unique(delivery_item)"`
	ShipmentQty int64 `json:"shipmentQty"`
	ReceiptQty  int64 `json:"receiptQty"`
	Committed   `xorm:"extends"`
}

func (Delivery) Create(ctx context.Context, param DeliveryCreateDto) (Delivery, error) {
	var items []DeliveryItem
	for _, i := range param.Items {
		item := DeliveryItem{
			SkuId:       i.SkuId,
			ShipmentQty: i.Qty,
			ReceiptQty:  0,
			Committed:   Committed{}.newCommitted(param.CreatedBy),
		}
		items = append(items, item)
	}
	d := Delivery{
		ShipmentLocationId: param.ShipmentLocationId,
		ReceiptLocationId:  param.ReceiptLocationId,
		WaybillNo:          param.WaybillNo,
		BoxNo:              param.BoxNo,
		Status:             Shipment.Code,
		Items:              items,
		Committed:          Committed{}.newCommitted(param.CreatedBy),
	}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryTableName).
		Insert(&d); err != nil {
		return Delivery{}, err
	}
	for i := range d.Items {
		d.Items[i].DeliveryId = d.Id
		d.Items[i].Committed = d.Committed
	}
	counts, err := DeliveryItem{}.insertDeliveryItems(ctx, d.Items)
	if err != nil {
		return Delivery{}, err
	}
	if counts != int64(len(d.Items)) {
		return Delivery{}, errors.New("fail to insert partial data")
	}
	return d, nil
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

func (Delivery) getById(ctx context.Context, id int64) (Delivery, error) {
	delivery := Delivery{}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryTableName).Alias("d").
		Where("d.id = ?", id).
		Get(&delivery); err != nil {
		return Delivery{}, err
	}
	if err := factory.
		DB(ctx).
		Table(DeliveryItemTableName).Alias("di").
		Where("di.delivery_id = ?", delivery.Id).
		Find(&delivery.Items); err != nil {
		return Delivery{}, err
	}
	return delivery, nil
}

func (Delivery) Receipt(ctx context.Context, param DeliveryReceiptDto) (Delivery, error) {
	for _, item := range param.Items {
		di := DeliveryItem{
			ReceiptQty: item.Qty,
			Committed: Committed{
				UpdatedBy: param.UpdatedBy,
			},
		}
		if _, err := factory.
			DB(ctx).
			Table(DeliveryItemTableName).
			Where("delivery_id = ? and sku_id = ?", param.DeliveryId, item.SkuId).
			Cols("receipt_qty, updated_by").
			Update(di); err != nil {
			return Delivery{}, err
		}
	}
	d := Delivery{
		Id:     param.DeliveryId,
		Status: Receipt.Code,
		Committed: Committed{
			UpdatedBy: param.UpdatedBy,
		},
	}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryTableName).
		ID(d.Id).
		Cols("status, updated_by").
		Update(&d); err != nil {
		return Delivery{}, err
	}
	return d, nil
}
