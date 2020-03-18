package controllers

import (
	"net/http"

	"github.com/hublabs/ehub-delivery-api/models"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type DeliveryController struct{}

func (c DeliveryController) Init(g echoswagger.ApiGroup) {
	g.POST("", c.Create).
		AddParamBody([]models.DeliveryCreateDto{}, "[]DeliveryCreateDto", "[]DeliveryCreateDto", true)
}

func (DeliveryController) Create(c echo.Context) error {
	var params []models.DeliveryCreateDto
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	for _, param := range params {
		d, err := param.Translate()
		if err != err {
			return c.JSON(http.StatusBadRequest, err)
		}
		if err := d.Create(c.Request().Context()); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return ReturnApiSucc(c, http.StatusCreated, int64(len(params)), nil)
}
