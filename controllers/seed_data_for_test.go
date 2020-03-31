package controllers

import (
	"github.com/go-xorm/xorm"
	"github.com/hublabs/delivery-api/models"
)

var (
	intQty    = int64(10)
	committed = models.Committed{
		CreatedBy: "li.dongxun",
	}
	SeedStocksForStore = []models.StockForStore{
		{Id: 1, LocationId: 100, SkuId: 10000, Qty: intQty, Committed: committed},
		{Id: 2, LocationId: 100, SkuId: 10001, Qty: intQty, Committed: committed},
		{Id: 3, LocationId: 100, SkuId: 10002, Qty: intQty, Committed: committed},
		{Id: 4, LocationId: 100, SkuId: 10003, Qty: intQty, Committed: committed},
		{Id: 5, LocationId: 200, SkuId: 20000, Qty: intQty, Committed: committed},
		{Id: 6, LocationId: 200, SkuId: 20001, Qty: intQty, Committed: committed},
		{Id: 7, LocationId: 200, SkuId: 20002, Qty: intQty, Committed: committed},
		{Id: 8, LocationId: 200, SkuId: 20003, Qty: intQty, Committed: committed},
		{Id: 9, LocationId: 300, SkuId: 30000, Qty: intQty, Committed: committed},
		{Id: 10, LocationId: 300, SkuId: 30001, Qty: intQty, Committed: committed},
		{Id: 11, LocationId: 300, SkuId: 30002, Qty: intQty, Committed: committed},
		{Id: 12, LocationId: 300, SkuId: 30003, Qty: intQty, Committed: committed},
	}
	SeedStocksForPlant = []models.StockForPlant{
		{Id: 1, LocationId: 10, SkuId: 11000, Qty: intQty, Committed: committed},
		{Id: 2, LocationId: 10, SkuId: 11001, Qty: intQty, Committed: committed},
		{Id: 3, LocationId: 10, SkuId: 11002, Qty: intQty, Committed: committed},
		{Id: 4, LocationId: 10, SkuId: 11003, Qty: intQty, Committed: committed},
		{Id: 5, LocationId: 20, SkuId: 21000, Qty: intQty, Committed: committed},
		{Id: 6, LocationId: 20, SkuId: 21001, Qty: intQty, Committed: committed},
		{Id: 7, LocationId: 20, SkuId: 21002, Qty: intQty, Committed: committed},
		{Id: 8, LocationId: 20, SkuId: 21003, Qty: intQty, Committed: committed},
		{Id: 9, LocationId: 30, SkuId: 31000, Qty: intQty, Committed: committed},
		{Id: 10, LocationId: 30, SkuId: 31001, Qty: intQty, Committed: committed},
		{Id: 11, LocationId: 30, SkuId: 31002, Qty: intQty, Committed: committed},
		{Id: 12, LocationId: 30, SkuId: 31003, Qty: intQty, Committed: committed},
	}
)

func CreateSeedData(xormEngine *xorm.Engine) error {
	if _, err := xormEngine.Table(models.StockForStoreTableName).Insert(&SeedStocksForStore); err != nil {
		return err
	}
	if _, err := xormEngine.Table(models.StockForPlantTableName).Insert(&SeedStocksForPlant); err != nil {
		return err
	}
	return nil
}
