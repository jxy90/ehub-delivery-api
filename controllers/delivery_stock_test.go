package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hublabs/delivery-api/models"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)

func TestDeliveryStock(t *testing.T) {
	const (
		TestShipAction    = "Ship"
		TestReceiveAction = "Receive"

		shippedBy = "shipper"
		receiptBy = "receiver"
	)

	t.Run(models.Store.Code, func(t *testing.T) {
		items1 := []map[string]interface{}{
			{
				"skuId": SeedStocksForStore[0].SkuId,
				"qty":   10,
			},
			{
				"skuId": SeedStocksForStore[1].SkuId,
				"qty":   2,
			},
		}
		items2 := []map[string]interface{}{
			{
				"skuId": SeedStocksForStore[2].SkuId,
				"qty":   10,
			},
			{
				"skuId": SeedStocksForStore[3].SkuId,
				"qty":   2,
			},
		}
		var (
			param1 = map[string]interface{}{
				"shipmentLocationId": SeedStocksForStore[0].LocationId,
				"receiptLocationId":  SeedStocksForStore[5].LocationId,
				"waybillNo":          "100",
				"boxNo":              "100",
				"type":               models.Store.Code,
				"items":              items1,
				"shippedBy":          shippedBy,
				"receiptBy":          receiptBy,
			}
			param2 = map[string]interface{}{
				"shipmentLocationId": SeedStocksForStore[0].LocationId,
				"receiptLocationId":  SeedStocksForStore[5].LocationId,
				"waybillNo":          "200",
				"boxNo":              "200",
				"type":               models.Store.Code,
				"items":              items2,
				"shippedBy":          shippedBy,
				"receiptBy":          receiptBy,
			}
		)

		t.Run(TestShipAction, func(t *testing.T) {
			shipmentBody := []map[string]interface{}{
				{
					"shipmentLocationId": param1["shipmentLocationId"],
					"receiptLocationId":  param1["receiptLocationId"],
					"waybillNo":          param1["waybillNo"],
					"boxNo":              param1["boxNo"],
					"type":               param1["type"],
					"items":              param1["items"],
					"shippedBy":          param1["shippedBy"],
				}, {
					"shipmentLocationId": param2["shipmentLocationId"],
					"receiptLocationId":  param2["receiptLocationId"],
					"waybillNo":          param2["waybillNo"],
					"boxNo":              param2["boxNo"],
					"type":               param2["type"],
					"items":              param2["items"],
					"shippedBy":          param2["shippedBy"],
				},
			}
			marshaledShipmentBody, _ := json.Marshal(shipmentBody)
			req := httptest.NewRequest(echo.POST, "/v1/delivery", bytes.NewReader(marshaledShipmentBody))
			setHeader(req)
			rec := httptest.NewRecorder()

			// api status check
			test.Ok(t, handleWithFilter(DeliveryController{}.Ship, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)

			// stock check
			afterShipStock0 := models.StockForStore{
				LocationId: SeedStocksForStore[0].LocationId,
				SkuId:      SeedStocksForStore[0].SkuId,
			}
			afterShipStock1 := models.StockForStore{
				LocationId: SeedStocksForStore[0].LocationId,
				SkuId:      SeedStocksForStore[1].SkuId,
			}
			afterShipStock2 := models.StockForStore{
				LocationId: SeedStocksForStore[0].LocationId,
				SkuId:      SeedStocksForStore[2].SkuId,
			}
			afterShipStock3 := models.StockForStore{
				LocationId: SeedStocksForStore[0].LocationId,
				SkuId:      SeedStocksForStore[3].SkuId,
			}

			foundStock0, _ := afterShipStock0.GetOne(dbContext)
			foundStock1, _ := afterShipStock1.GetOne(dbContext)
			foundStock2, _ := afterShipStock2.GetOne(dbContext)
			foundStock3, _ := afterShipStock3.GetOne(dbContext)

			test.Equals(t, foundStock0.Qty, int64(0))
			test.Equals(t, foundStock1.Qty, int64(8))
			test.Equals(t, foundStock2.Qty, int64(0))
			test.Equals(t, foundStock3.Qty, int64(8))
		})

		t.Run(TestReceiveAction, func(t *testing.T) {
			receiptBody := []map[string]interface{}{
				{
					"deliveryId":        1,
					"receiptLocationId": param1["receiptLocationId"],
					"type":              param1["type"],
					"items":             param1["items"],
					"receiptedBy":       param1["receiptBy"],
				},
				{
					"deliveryId":        2,
					"receiptLocationId": param2["receiptLocationId"],
					"type":              param2["type"],
					"items":             param2["items"],
					"receiptedBy":       param2["receiptBy"],
				},
			}
			marshaledReceiptBody, _ := json.Marshal(receiptBody)
			req := httptest.NewRequest(echo.PATCH, "/v1/delivery", bytes.NewReader(marshaledReceiptBody))
			setHeader(req)
			rec := httptest.NewRecorder()

			// api status check
			test.Ok(t, handleWithFilter(DeliveryController{}.Receive, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusAccepted, rec.Code)

			// stock check
			afterReceiveStock0 := models.StockForStore{
				LocationId: SeedStocksForStore[5].LocationId,
				SkuId:      SeedStocksForStore[0].SkuId,
			}
			afterReceiveStock1 := models.StockForStore{
				LocationId: SeedStocksForStore[5].LocationId,
				SkuId:      SeedStocksForStore[1].SkuId,
			}
			afterReceiveStock2 := models.StockForStore{
				LocationId: SeedStocksForStore[5].LocationId,
				SkuId:      SeedStocksForStore[2].SkuId,
			}
			afterReceiveStock3 := models.StockForStore{
				LocationId: SeedStocksForStore[5].LocationId,
				SkuId:      SeedStocksForStore[3].SkuId,
			}

			foundStock0, _ := afterReceiveStock0.GetOne(dbContext)
			foundStock1, _ := afterReceiveStock1.GetOne(dbContext)
			foundStock2, _ := afterReceiveStock2.GetOne(dbContext)
			foundStock3, _ := afterReceiveStock3.GetOne(dbContext)

			test.Equals(t, foundStock0.Qty, int64(10))
			test.Equals(t, foundStock1.Qty, int64(2))
			test.Equals(t, foundStock2.Qty, int64(10))
			test.Equals(t, foundStock3.Qty, int64(2))
		})
	})

	t.Run(models.Plant.Code, func(t *testing.T) {
		items1 := []map[string]interface{}{
			{
				"skuId": SeedStocksForPlant[0].SkuId,
				"qty":   10,
			},
			{
				"skuId": SeedStocksForPlant[1].SkuId,
				"qty":   2,
			},
		}
		items2 := []map[string]interface{}{
			{
				"skuId": SeedStocksForPlant[2].SkuId,
				"qty":   10,
			},
			{
				"skuId": SeedStocksForPlant[3].SkuId,
				"qty":   2,
			},
		}
		var (
			param1 = map[string]interface{}{
				"shipmentLocationId": SeedStocksForPlant[0].LocationId,
				"receiptLocationId":  SeedStocksForPlant[5].LocationId,
				"waybillNo":          "100",
				"boxNo":              "100",
				"type":               models.Plant.Code,
				"items":              items1,
				"shippedBy":          shippedBy,
				"receiptBy":          receiptBy,
			}
			param2 = map[string]interface{}{
				"shipmentLocationId": SeedStocksForPlant[0].LocationId,
				"receiptLocationId":  SeedStocksForPlant[5].LocationId,
				"waybillNo":          "200",
				"boxNo":              "200",
				"type":               models.Plant.Code,
				"items":              items2,
				"shippedBy":          shippedBy,
				"receiptBy":          receiptBy,
			}
		)

		t.Run(TestShipAction, func(t *testing.T) {
			shipmentBody := []map[string]interface{}{
				{
					"shipmentLocationId": param1["shipmentLocationId"],
					"receiptLocationId":  param1["receiptLocationId"],
					"waybillNo":          param1["waybillNo"],
					"boxNo":              param1["boxNo"],
					"type":               param1["type"],
					"items":              param1["items"],
					"shippedBy":          param1["shippedBy"],
				}, {
					"shipmentLocationId": param2["shipmentLocationId"],
					"receiptLocationId":  param2["receiptLocationId"],
					"waybillNo":          param2["waybillNo"],
					"boxNo":              param2["boxNo"],
					"type":               param2["type"],
					"items":              param2["items"],
					"shippedBy":          param2["shippedBy"],
				},
			}
			marshaledShipmentBody, _ := json.Marshal(shipmentBody)
			req := httptest.NewRequest(echo.POST, "/v1/delivery", bytes.NewReader(marshaledShipmentBody))
			setHeader(req)
			rec := httptest.NewRecorder()

			// api status check
			test.Ok(t, handleWithFilter(DeliveryController{}.Ship, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)

			// stock check
			afterShipStock0 := models.StockForPlant{
				LocationId: SeedStocksForPlant[0].LocationId,
				SkuId:      SeedStocksForPlant[0].SkuId,
			}
			afterShipStock1 := models.StockForPlant{
				LocationId: SeedStocksForPlant[0].LocationId,
				SkuId:      SeedStocksForPlant[1].SkuId,
			}
			afterShipStock2 := models.StockForPlant{
				LocationId: SeedStocksForPlant[0].LocationId,
				SkuId:      SeedStocksForPlant[2].SkuId,
			}
			afterShipStock3 := models.StockForPlant{
				LocationId: SeedStocksForPlant[0].LocationId,
				SkuId:      SeedStocksForPlant[3].SkuId,
			}

			foundStock0, _ := afterShipStock0.GetOne(dbContext)
			foundStock1, _ := afterShipStock1.GetOne(dbContext)
			foundStock2, _ := afterShipStock2.GetOne(dbContext)
			foundStock3, _ := afterShipStock3.GetOne(dbContext)

			test.Equals(t, foundStock0.Qty, int64(0))
			test.Equals(t, foundStock1.Qty, int64(8))
			test.Equals(t, foundStock2.Qty, int64(0))
			test.Equals(t, foundStock3.Qty, int64(8))
		})

		t.Run(TestReceiveAction, func(t *testing.T) {
			receiptBody := []map[string]interface{}{
				{
					"deliveryId":        1,
					"receiptLocationId": param1["receiptLocationId"],
					"type":              param1["type"],
					"items":             param1["items"],
					"receiptedBy":       param1["receiptBy"],
				},
				{
					"deliveryId":        2,
					"receiptLocationId": param2["receiptLocationId"],
					"type":              param2["type"],
					"items":             param2["items"],
					"receiptedBy":       param2["receiptBy"],
				},
			}
			marshaledReceiptBody, _ := json.Marshal(receiptBody)
			req := httptest.NewRequest(echo.PATCH, "/v1/delivery", bytes.NewReader(marshaledReceiptBody))
			setHeader(req)
			rec := httptest.NewRecorder()

			// api status check
			test.Ok(t, handleWithFilter(DeliveryController{}.Receive, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusAccepted, rec.Code)

			// stock check
			afterReceiveStock0 := models.StockForPlant{
				LocationId: SeedStocksForPlant[5].LocationId,
				SkuId:      SeedStocksForPlant[0].SkuId,
			}
			afterReceiveStock1 := models.StockForPlant{
				LocationId: SeedStocksForPlant[5].LocationId,
				SkuId:      SeedStocksForPlant[1].SkuId,
			}
			afterReceiveStock2 := models.StockForPlant{
				LocationId: SeedStocksForPlant[5].LocationId,
				SkuId:      SeedStocksForPlant[2].SkuId,
			}
			afterReceiveStock3 := models.StockForPlant{
				LocationId: SeedStocksForPlant[5].LocationId,
				SkuId:      SeedStocksForPlant[3].SkuId,
			}

			foundStock0, _ := afterReceiveStock0.GetOne(dbContext)
			foundStock1, _ := afterReceiveStock1.GetOne(dbContext)
			foundStock2, _ := afterReceiveStock2.GetOne(dbContext)
			foundStock3, _ := afterReceiveStock3.GetOne(dbContext)

			test.Equals(t, foundStock0.Qty, int64(10))
			test.Equals(t, foundStock1.Qty, int64(2))
			test.Equals(t, foundStock2.Qty, int64(10))
			test.Equals(t, foundStock3.Qty, int64(2))
		})
	})
}
