package models

import (
	"context"
)

var (
	DeliveryForStoreTableName     = "dlv_delivery_for_store"
	DeliveryItemForStoreTableName = "dlv_delivery_item_for_store"
	DeliveryForPlantTableName     = "dlv_delivery_for_plant"
	DeliveryItemForPlantTableName = "dlv_delivery_item_for_plant"
)

func (d *DeliveryForStore) TableName() string {
	return DeliveryForStoreTableName
}

func (di *DeliveryItemForStore) TableName() string {
	return DeliveryItemForStoreTableName
}

func (d *DeliveryForPlant) TableName() string {
	return DeliveryForPlantTableName
}

func (di *DeliveryItemForPlant) TableName() string {
	return DeliveryItemForPlantTableName
}

type DeliveryProcessor interface {
	ship(ctx context.Context) error
	receive(ctx context.Context) error
	calculateStock(ctx context.Context) error
	translateToStockByStatus() ([]StockProcessor, error)
}

type DeliveryTranslator interface {
	translateToDeliveryByType() (DeliveryProcessor, error)
}

func Ship(ctx context.Context, shipmentDto ShipmentDto) (DeliveryProcessor, error) {
	d, err := shipmentDto.translateToDeliveryByType()
	if err != nil {
		return nil, err
	}
	if err := d.ship(ctx); err != nil {
		return nil, err
	}
	if err := d.calculateStock(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

func Receive(ctx context.Context, receiptDto ReceiptDto) (DeliveryProcessor, error) {
	d, err := receiptDto.translateToDeliveryByType()
	if err != nil {
		return nil, err
	}
	if err := d.receive(ctx); err != nil {
		return nil, err
	}
	if err := d.calculateStock(ctx); err != nil {
		return nil, err
	}
	return d, nil
}
