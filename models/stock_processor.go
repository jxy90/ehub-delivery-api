package models

import (
	"context"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize"
)

var (
	StockForStoreTableName = "stk_stock_for_store"
	StockForPlantTableName = "stk_stock_for_plant"
)

func (s *StockForStore) TableName() string {
	return StockForStoreTableName
}

func (s *StockForPlant) TableName() string {
	return StockForPlantTableName
}

type StockProcessor interface {
	bulkCreateStockFromDto(ctx context.Context, param StockCreateDto) (int64, error)
	bulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy string, excel *excelize.File) (int64, error)
	upsertQty(ctx context.Context) error
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

func BulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy, deliveryType string, excel *excelize.File) (int64, error) {
	if deliveryType == Store.Code {
		count, err := StockForStore{}.bulkCreateStockFromExcel(ctx, locationId, createdBy, excel)
		if err != nil {
			return 0, err
		}

		return count, nil
	}
	if deliveryType == Plant.Code {
		count, err := StockForPlant{}.bulkCreateStockFromExcel(ctx, locationId, createdBy, excel)
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	return 0, errors.New("not supported type, supported types are [store, plant]")
}
