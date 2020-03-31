package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hublabs/common/api"
	"github.com/hublabs/delivery-api/models"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type StockController struct{}

func (c StockController) Init(g echoswagger.ApiGroup) {
	g.POST("", c.Create).
		AddParamBody(models.StockCreateDto{}, "StockCreateDto", "StockCreateDto", false)
}

func (StockController) Create(c echo.Context) error {
	if strings.Index(c.Request().Header.Get("Content-Type"), "application/json") == 0 {
		var param models.StockCreateDto
		if err := c.Bind(&param); err != nil {
			return ReturnError(c, http.StatusBadRequest, api.Error{
				Message: err.Error(),
			})
		}
		if err := c.Validate(param); err != nil {
			return ReturnError(c, http.StatusBadRequest, api.Error{
				Message: err.Error(),
			})
		}
		count, err := models.BulkCreateStockFromDto(c.Request().Context(), param)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}
		return ReturnSuccessWithTotalCountAndItems(c, http.StatusOK, count, nil)
	}

	if strings.Index(c.Request().Header.Get("Content-Type"), "multipart/form-data") == 0 {
		// ref: https://echo.labstack.com/cookbook/file-upload
		location := c.FormValue("locationId")
		locationId, err := strconv.ParseInt(location, 10, 64)
		if err != nil {
			var errorParams []string
			errorParams = append(errorParams, location)
			return ReturnErrorWithParameter(c, errorParams)
		}
		createdBy := c.FormValue("createdBy")
		file, err := c.FormFile("stockFile")
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}

		excel, err := excelize.OpenFile(file.Filename)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}

		stockType := c.FormValue("type")
		count, err := models.BulkCreateStockFromExcel(c.Request().Context(), locationId, createdBy, stockType, excel)
		if err != nil {
			return ReturnError(c, http.StatusInternalServerError, api.Error{
				Message: err.Error(),
			})
		}
		return ReturnSuccessWithTotalCountAndItems(c, http.StatusOK, count, nil)
	}

	return ReturnError(c, http.StatusInternalServerError, api.Error{
		Message: "not supported Content-Type, supported Content-Types are [application/json, multipart/form-data]",
	})
}
