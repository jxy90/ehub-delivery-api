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
	g.POST("", c.Create).
		AddParamBody([]models.DeliveryCreateDto{}, "[]DeliveryCreateDto", "[]DeliveryCreateDto", true)
	g.PATCH("", c.Receipt).
		AddParamBody([]models.DeliveryReceiptDto{}, "[]DeliveryReceiptDto", "[]DeliveryReceiptDto", true)
}

func (DeliveryController) Create(c echo.Context) error {
	var params []models.DeliveryCreateDto
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var created []models.Delivery
	for _, param := range params {
		d, err := models.Delivery{}.Create(c.Request().Context(), param)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}
		created = append(created, d)
	}

	return ReturnSuccessWithTotalCountAndItems(c, http.StatusCreated, int64(len(params)), created)
}

func (DeliveryController) Receipt(c echo.Context) error {
	var params []models.DeliveryReceiptDto
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var updateds []models.Delivery
	for _, param := range params {
		d, err := models.Delivery{}.Receipt(c.Request().Context(), param)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}
		updateds = append(updateds, d)
	}

	return ReturnSuccessWithTotalCountAndItems(c, http.StatusCreated, int64(len(updateds)), updateds)
}
