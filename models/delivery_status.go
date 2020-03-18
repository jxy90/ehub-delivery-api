package models

type DeliveryStatus struct {
	Code                string `json:"code"`
	ShipmentStockChange bool   `json:"shipmentStockChange"`
	ReceiptStockChange  bool   `json:"receiptStockChange"`
	Description         string `json:"description"`
}

var (
	Shipment = DeliveryStatus{
		Code:                "S",
		ShipmentStockChange: true,
		ReceiptStockChange:  false,
		Description:         "change shipment location stock",
	}
	Receipt = DeliveryStatus{
		Code:                "R",
		ShipmentStockChange: false,
		ReceiptStockChange:  true,
		Description:         "change receipt location stock",
	}
)
