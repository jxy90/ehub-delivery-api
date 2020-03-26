package models

import (
	"context"
	"errors"
	"github.com/hublabs/delivery-api/factory"
)

type DeliveryForStore struct {
	Id                 int64                  `json:"id"`
	ShipmentLocationId int64                  `json:"shipmentLocationId" xorm:"index unique(delivery)"`
	ReceiptLocationId  int64                  `json:"receiptLocationId" xorm:"index unique(delivery)"`
	WaybillNo          string                 `json:"waybillNo" xorm:"index unique(delivery)"`
	BoxNo              string                 `json:"boxNo" xorm:"index unique(delivery)"`
	Status             string                 `json:"status"`
	Items              []DeliveryItemForStore `json:"items" xorm:"-"`
	Committed          `xorm:"extends"`
}

type DeliveryItemForStore struct {
	Id          int64 `json:"id"`
	DeliveryId  int64 `json:"deliveryId" xorm:"index unique(delivery_item)"`
	SkuId       int64 `json:"skuId" xorm:"index unique(delivery_item)"`
	ShipmentQty int64 `json:"shipmentQty"`
	ReceiptQty  int64 `json:"receiptQty"`
	Committed   `xorm:"extends"`
}

func (d *DeliveryForStore) ship(ctx context.Context) error {
	if _, err := factory.
		DB(ctx).
		Table(DeliveryForStoreTableName).
		Insert(d); err != nil {
		return err
	}
	for i := range d.Items {
		d.Items[i].DeliveryId = d.Id
		d.Items[i].Committed = d.Committed
	}
	counts, err := DeliveryItemForStore{}.insertDeliveryItems(ctx, d.Items)
	if err != nil {
		return err
	}
	if counts != int64(len(d.Items)) {
		return errors.New("fail to insert partial data")
	}
	return nil
}

func (DeliveryItemForStore) insertDeliveryItems(ctx context.Context, items []DeliveryItemForStore) (int64, error) {
	counts, err := factory.
		DB(ctx).
		Table(DeliveryItemForStoreTableName).
		Insert(items)
	if err != nil {
		return 0, err
	}
	return counts, nil
}

func (DeliveryItemForStore) getById(ctx context.Context, id int64) (DeliveryForStore, error) {
	delivery := DeliveryForStore{}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryForStoreTableName).Alias("d").
		Where("d.id = ?", id).
		Get(&delivery); err != nil {
		return DeliveryForStore{}, err
	}
	if err := factory.
		DB(ctx).
		Table(DeliveryItemForStoreTableName).Alias("di").
		Where("di.delivery_id = ?", delivery.Id).
		Find(&delivery.Items); err != nil {
		return DeliveryForStore{}, err
	}
	return delivery, nil
}

func (d *DeliveryForStore) receive(ctx context.Context) error {
	for _, item := range d.Items {
		if _, err := factory.
			DB(ctx).
			Table(DeliveryItemForStoreTableName).
			Where("delivery_id = ? and sku_id = ?", item.DeliveryId, item.SkuId).
			Cols("receipt_qty, updated_by").
			Update(&item); err != nil {
			return err
		}
	}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryForStoreTableName).
		ID(d.Id).
		Cols("status, updated_by").
		Update(d); err != nil {
		return err
	}
	return nil
}

func (d *DeliveryForStore) calculateStock(ctx context.Context) error {
	stocks, err := d.translateToStockByStatus()
	if err != nil {
		return err
	}
	for _, stock := range stocks {
		inputStock := stock.(*StockForStore)
		foundStock, err := inputStock.GetOne(ctx)
		if err != nil {
			return err
		}
		updateStock := StockForStore{
			LocationId: inputStock.LocationId,
			SkuId:      inputStock.SkuId,
			Qty:        inputStock.Qty + foundStock.Qty,
			Committed:  Committed{UpdatedBy: inputStock.UpdatedBy},
		}
		if err := updateStock.upsertQty(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeliveryForStore) translateToStockByStatus() ([]StockProcessor, error) {
	var stocks []StockProcessor
	switch d.Status {
	case Shipment.Code:
		for _, item := range d.Items {
			stock := &StockForStore{
				LocationId: d.ShipmentLocationId,
				SkuId:      item.SkuId,
				Qty:        item.ShipmentQty * -1,
				Committed: Committed{
					UpdatedBy: d.UpdatedBy,
				},
			}
			stocks = append(stocks, stock)
		}
		return stocks, nil
	case Receipt.Code:
		for _, item := range d.Items {
			stock := &StockForStore{
				LocationId: d.ReceiptLocationId,
				SkuId:      item.SkuId,
				Qty:        item.ReceiptQty,
				Committed: Committed{
					UpdatedBy: d.UpdatedBy,
				},
			}
			stocks = append(stocks, stock)
		}
		return stocks, nil
	}
	return nil, errors.New("not support status")
}

type DeliveryForPlant struct {
	Id                 int64                  `json:"id"`
	ShipmentLocationId int64                  `json:"shipmentLocationId" xorm:"index unique(delivery)"`
	ReceiptLocationId  int64                  `json:"receiptLocationId" xorm:"index unique(delivery)"`
	WaybillNo          string                 `json:"waybillNo" xorm:"index unique(delivery)"`
	BoxNo              string                 `json:"boxNo" xorm:"index unique(delivery)"`
	Status             string                 `json:"status"`
	Items              []DeliveryItemForPlant `json:"items" xorm:"-"`
	PlatformOrderId    string                 `json:"platformOrderId" xorm:"index"`
	Committed          `xorm:"extends"`
}

type DeliveryItemForPlant struct {
	Id          int64 `json:"id"`
	DeliveryId  int64 `json:"deliveryId" xorm:"index unique(delivery_item)"`
	SkuId       int64 `json:"skuId" xorm:"index unique(delivery_item)"`
	ShipmentQty int64 `json:"shipmentQty"`
	ReceiptQty  int64 `json:"receiptQty"`
	Committed   `xorm:"extends"`
}

func (d *DeliveryForPlant) ship(ctx context.Context) error {
	if _, err := factory.
		DB(ctx).
		Table(DeliveryForPlantTableName).
		Insert(d); err != nil {
		return err
	}
	for i := range d.Items {
		d.Items[i].DeliveryId = d.Id
		d.Items[i].Committed = d.Committed
	}
	counts, err := DeliveryItemForPlant{}.insertDeliveryItems(ctx, d.Items)
	if err != nil {
		return err
	}
	if counts != int64(len(d.Items)) {
		return errors.New("fail to insert partial data")
	}
	return nil
}

func (DeliveryItemForPlant) insertDeliveryItems(ctx context.Context, items []DeliveryItemForPlant) (int64, error) {
	counts, err := factory.
		DB(ctx).
		Table(DeliveryItemForPlantTableName).
		Insert(items)
	if err != nil {
		return 0, err
	}
	return counts, nil
}

func (DeliveryForPlant) getById(ctx context.Context, id int64) (DeliveryForPlant, error) {
	delivery := DeliveryForPlant{}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryForPlantTableName).Alias("d").
		Where("d.id = ?", id).
		Get(&delivery); err != nil {
		return DeliveryForPlant{}, err
	}
	if err := factory.
		DB(ctx).
		Table(DeliveryItemForPlantTableName).Alias("di").
		Where("di.delivery_id = ?", delivery.Id).
		Find(&delivery.Items); err != nil {
		return DeliveryForPlant{}, err
	}
	return delivery, nil
}

func (d *DeliveryForPlant) receive(ctx context.Context) error {
	for _, item := range d.Items {
		if _, err := factory.
			DB(ctx).
			Table(DeliveryItemForPlantTableName).
			Where("delivery_id = ? and sku_id = ?", item.DeliveryId, item.SkuId).
			Cols("receipt_qty, updated_by").
			Update(&item); err != nil {
			return err
		}
	}
	if _, err := factory.
		DB(ctx).
		Table(DeliveryForPlantTableName).
		ID(d.Id).
		Cols("status, updated_by").
		Update(d); err != nil {
		return err
	}
	return nil
}

func (d *DeliveryForPlant) calculateStock(ctx context.Context) error {
	stocks, err := d.translateToStockByStatus()
	if err != nil {
		return err
	}
	for _, stock := range stocks {
		inputStock := stock.(*StockForPlant)
		foundStock, err := inputStock.GetOne(ctx)
		if err != nil {
			return err
		}
		updateStock := StockForPlant{
			LocationId: inputStock.LocationId,
			SkuId:      inputStock.SkuId,
			Qty:        inputStock.Qty + foundStock.Qty,
			Committed:  Committed{UpdatedBy: inputStock.UpdatedBy},
		}
		if err := updateStock.upsertQty(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *DeliveryForPlant) translateToStockByStatus() ([]StockProcessor, error) {
	var stocks []StockProcessor
	switch d.Status {
	case Shipment.Code:
		for _, item := range d.Items {
			stock := &StockForPlant{
				LocationId: d.ShipmentLocationId,
				SkuId:      item.SkuId,
				Qty:        item.ShipmentQty * -1,
				Committed: Committed{
					UpdatedBy: d.UpdatedBy,
				},
			}
			stocks = append(stocks, stock)
		}
		return stocks, nil
	case Receipt.Code:
		for _, item := range d.Items {
			stock := &StockForPlant{
				LocationId: d.ReceiptLocationId,
				SkuId:      item.SkuId,
				Qty:        item.ReceiptQty,
				Committed: Committed{
					UpdatedBy: d.UpdatedBy,
				},
			}
			stocks = append(stocks, stock)
		}
		return stocks, nil
	}
	return nil, errors.New("not support status")
}
