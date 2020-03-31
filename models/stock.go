package models

import (
	"context"
	"strconv"

	"github.com/hublabs/delivery-api/factory"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type StockForStore struct {
	Id         int64 `json:"id"`
	LocationId int64 `json:"locationId" xorm:"index unique(stock)"`
	SkuId      int64 `json:"skuId" xorm:"index unique(stock)"`
	Qty        int64 `json:"qty"`
	Committed  `xorm:"extends"`
}

func (StockForStore) bulkCreateStockFromDto(ctx context.Context, param StockCreateDto) (int64, error) {
	var stocks []StockForStore
	for _, item := range param.Items {
		stock := StockForStore{
			LocationId: param.LocationId,
			SkuId:      item.SkuId,
			Qty:        item.Qty,
			Committed:  Committed{}.newCommitted(param.CreatedBy),
		}
		stocks = append(stocks, stock)
	}
	if _, err := factory.
		DB(ctx).
		Table(StockForStoreTableName).
		Insert(&stocks); err != nil {
		return 0, err
	}
	return int64(len(stocks)), nil
}

func (StockForStore) bulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy string, excel *excelize.File) (int64, error) {
	rows := excel.GetRows("Sheet1")
	cellMaps := make([]map[string]int64, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var cellMap map[string]int64
		for i, colCell := range row {
			if i%2 == 0 {
				cellMap = make(map[string]int64, 0)
				skuId, err := strconv.ParseInt(colCell, 10, 64)
				if err != nil {
					return 0, err
				}
				cellMap["skuId"] = skuId
			} else {
				qty, err := strconv.ParseInt(colCell, 10, 64)
				if err != nil {
					return 0, err
				}
				if cellMap != nil {
					cellMap["qty"] = qty
					cellMaps = append(cellMaps, cellMap)
				}
			}
		}
	}
	var stocks []StockForStore
	for _, val := range cellMaps {
		stock := StockForStore{
			LocationId: locationId,
			SkuId:      val["skuId"],
			Qty:        val["qty"],
			Committed:  Committed{}.newCommitted(createdBy),
		}
		stocks = append(stocks, stock)
	}
	if _, err := factory.
		DB(ctx).
		Table(StockForStoreTableName).
		Insert(&stocks); err != nil {
		return 0, err
	}
	return int64(len(stocks)), nil
}

func (s *StockForStore) GetOne(ctx context.Context) (StockForStore, error) {
	var found StockForStore
	if _, err := factory.
		DB(ctx).
		Table(StockForStoreTableName).
		Where("location_id = ? and sku_id = ?", s.LocationId, s.SkuId).
		Get(&found); err != nil {
		return StockForStore{}, err
	}
	return found, nil
}

func (s *StockForStore) upsertQty(ctx context.Context) error {
	count, err := factory.
		DB(ctx).
		Table(StockForStoreTableName).
		Where("location_id = ? and sku_id = ?", s.LocationId, s.SkuId).
		Cols("qty, updated_by").
		Update(s)
	if err != nil {
		return err
	}
	if count == 0 {
		creating := StockForStore{
			LocationId: s.LocationId,
			SkuId:      s.SkuId,
			Qty:        s.Qty,
			Committed:  s.newCommitted(s.Committed.UpdatedBy),
		}
		if _, err := factory.
			DB(ctx).
			Table(StockForStoreTableName).
			Insert(creating); err != nil {
			return err
		}
	}
	return nil
}

type StockForPlant struct {
	Id         int64 `json:"id"`
	LocationId int64 `json:"locationId" xorm:"index unique(stock)"`
	SkuId      int64 `json:"skuId" xorm:"index unique(stock)"`
	Qty        int64 `json:"qty"`
	Committed  `xorm:"extends"`
}

func (StockForPlant) bulkCreateStockFromDto(ctx context.Context, param StockCreateDto) (int64, error) {
	var stocks []StockForPlant
	for _, item := range param.Items {
		stock := StockForPlant{
			LocationId: param.LocationId,
			SkuId:      item.SkuId,
			Qty:        item.Qty,
			Committed:  Committed{}.newCommitted(param.CreatedBy),
		}
		stocks = append(stocks, stock)
	}
	if _, err := factory.
		DB(ctx).
		Table(StockForPlantTableName).
		Insert(&stocks); err != nil {
		return 0, err
	}
	return int64(len(stocks)), nil
}

func (StockForPlant) bulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy string, excel *excelize.File) (int64, error) {
	rows := excel.GetRows("Sheet1")
	cellMaps := make([]map[string]int64, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		var cellMap map[string]int64
		for i, colCell := range row {
			if i%2 == 0 {
				cellMap = make(map[string]int64, 0)
				skuId, err := strconv.ParseInt(colCell, 10, 64)
				if err != nil {
					return 0, err
				}
				cellMap["skuId"] = skuId
			} else {
				qty, err := strconv.ParseInt(colCell, 10, 64)
				if err != nil {
					return 0, err
				}
				if cellMap != nil {
					cellMap["qty"] = qty
					cellMaps = append(cellMaps, cellMap)
				}
			}
		}
	}
	var stocks []StockForPlant
	for _, val := range cellMaps {
		stock := StockForPlant{
			LocationId: locationId,
			SkuId:      val["skuId"],
			Qty:        val["qty"],
			Committed:  Committed{}.newCommitted(createdBy),
		}
		stocks = append(stocks, stock)
	}
	if _, err := factory.
		DB(ctx).
		Table(StockForPlantTableName).
		Insert(&stocks); err != nil {
		return 0, err
	}
	return int64(len(stocks)), nil
}

func (s *StockForPlant) GetOne(ctx context.Context) (StockForPlant, error) {
	var found StockForPlant
	if _, err := factory.
		DB(ctx).
		Table(StockForPlantTableName).
		Where("location_id = ? and sku_id = ?", s.LocationId, s.SkuId).
		Get(&found); err != nil {
		return StockForPlant{}, err
	}
	return found, nil
}

func (s *StockForPlant) upsertQty(ctx context.Context) error {
	count, err := factory.
		DB(ctx).
		Table(StockForPlantTableName).
		Where("location_id = ? and sku_id = ?", s.LocationId, s.SkuId).
		Cols("qty, updated_by").
		Update(s)
	if err != nil {
		return err
	}
	if count == 0 {
		creating := StockForPlant{
			LocationId: s.LocationId,
			SkuId:      s.SkuId,
			Qty:        s.Qty,
			Committed:  s.newCommitted(s.Committed.UpdatedBy),
		}
		if _, err := factory.
			DB(ctx).
			Table(StockForPlantTableName).
			Insert(creating); err != nil {
			return err
		}
	}
	return nil
}
