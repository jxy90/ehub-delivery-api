package models

import (
	"context"
	"errors"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var (
	DeliveryForStoreTableName     = "dlv_delivery_for_store"
	DeliveryItemForStoreTableName = "dlv_delivery_item_for_store"
	DeliveryForPlantTableName     = "dlv_delivery_for_plant"
	DeliveryItemForPlantTableName = "dlv_delivery_item_for_plant"
	StockForStoreTableName        = "stk_stock_for_store"
	StockForPlantTableName        = "stk_stock_for_plant"
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

func (s *StockForStore) TableName() string {
	return StockForStoreTableName
}

func (s *StockForPlant) TableName() string {
	return StockForPlantTableName
}

type DeliveryProcessor interface {
	ship(ctx context.Context) error
	receive(ctx context.Context) error
}

func Ship(ctx context.Context, dto ShipmentDto) (DeliveryProcessor, error) {
	d, err := dto.translateToDelivery()
	if err != nil {
		return nil, err
	}
	if err := d.ship(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

func Receive(ctx context.Context, dto ReceiptDto) (DeliveryProcessor, error) {
	d, err := dto.translateToDelivery()
	if err != nil {
		return nil, err
	}
	if err := d.receive(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

type StockProcessor interface {
	bulkCreateStockFromDto(ctx context.Context, param StockCreateDto) (int64, error)
	bulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy string, excel *excelize.File) (int64, error)
}

func BulkCreateStockFromDto(ctx context.Context, param StockCreateDto) (int64, error) {
	if param.Type == Store.Code {
		count, err := StockForStore{}.bulkCreateStockFromDto(ctx, param)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	if param.Type == Plant.Code {
		count, err := StockForPlant{}.bulkCreateStockFromDto(ctx, param)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, errors.New("not supported type, supported types are [store, plant]")
}

func BulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy, stockType string, excel *excelize.File) (int64, error) {
	if stockType == Store.Code {
		count, err := StockForStore{}.bulkCreateStockFromExcel(ctx, locationId, createdBy, excel)
		if err != nil {
			return 0, err
		}

		return count, nil
	}
	if stockType == Plant.Code {
		count, err := StockForPlant{}.bulkCreateStockFromExcel(ctx, locationId, createdBy, excel)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, errors.New("not supported type, supported types are [store, plant]")
}
