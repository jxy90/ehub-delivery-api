package controllers

import (
	"net/http"

	"github.com/hublabs/common/api"
	"github.com/hublabs/ehub-delivery-api/models"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type DeliveryController struct{}

func (c DeliveryController) Init(g echoswagger.ApiGroup) {
	g.POST("", c.Ship).
		AddParamBody([]models.ShipmentDto{}, "[]ShipmentDto", "[]ShipmentDto", true)
	g.PATCH("", c.Receive).
		AddParamBody([]models.ReceiptDto{}, "[]ReceiptDto", "[]ReceiptDto", true)
}

func (DeliveryController) Ship(c echo.Context) error {
	var params []models.ShipmentDto
	if err := c.Bind(&params); err != nil {
		return ReturnError(c, http.StatusBadRequest, api.Error{
			Message: err.Error(),
		})
	}

	var created []models.DeliveryProcessor
	for _, param := range params {
		d, err := models.Ship(c.Request().Context(), param)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}
		created = append(created, d)
	}

	return ReturnSuccessWithTotalCountAndItems(c, http.StatusCreated, int64(len(params)), created)
}

func (DeliveryController) Receive(c echo.Context) error {
	var params []models.ReceiptDto
	if err := c.Bind(&params); err != nil {
		return ReturnError(c, http.StatusBadRequest, api.Error{
			Message: err.Error(),
		})
	}

	var updateds []models.DeliveryProcessor
	for _, param := range params {
		d, err := models.Receive(c.Request().Context(), param)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}
		updateds = append(updateds, d)
	}

	return ReturnSuccessWithTotalCountAndItems(c, http.StatusCreated, int64(len(updateds)), updateds)
}
