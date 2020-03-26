package models

import (
	"context"
	"strconv"

	"github.com/hublabs/ehub-delivery-api/factory"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var (
	StockTableName = "stk_stock"
)

func (s *Stock) TableName() string {
	return StockTableName
}

type Stock struct {
	Id         int64 `json:"id"`
	LocationId int64 `json:"locationId" xorm:"index unique(stock)"`
	SkuId      int64 `json:"skuId" xorm:"index unique(stock)"`
	Qty        int64 `json:"qty"`
	Committed  `xorm:"extends"`
}

func (Stock) BulkCreateStockFromDto(ctx context.Context, params []StockCreateDto) (int64, error) {
	var stocks []Stock
	for _, param := range params {
		for _, item := range param.Items {
			stock := Stock{
				LocationId: param.LocationId,
				SkuId:      item.SkuId,
				Qty:        item.Qty,
				Committed:  Committed{}.newCommitted(param.CreatedBy),
			}
			stocks = append(stocks, stock)
		}
	}
	if _, err := factory.
		DB(ctx).
		Table(StockTableName).
		Insert(&stocks); err != nil {
		return 0, err
	}
	return int64(len(stocks)), nil
}

func (Stock) BulkCreateStockFromExcel(ctx context.Context, locationId int64, createdBy string, excel *excelize.File) (int64, error) {
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
	var stocks []Stock
	for _, val := range cellMaps {
		stock := Stock{
			LocationId: locationId,
			SkuId:      val["skuId"],
			Qty:        val["qty"],
			Committed:  Committed{}.newCommitted(createdBy),
		}
		stocks = append(stocks, stock)
	}
	if _, err := factory.
		DB(ctx).
		Table(StockTableName).
		Insert(&stocks); err != nil {
		return 0, err
	}
	return int64(len(stocks)), nil
}
