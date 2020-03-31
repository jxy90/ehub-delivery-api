package controllers

import (
	"fmt"
	"net/http"

	"github.com/hublabs/common/api"

	"github.com/labstack/echo"
)

var ProjectName string = "[delivery-api]"

type SearchPageCount struct {
	SkipCount      int `query:"skipCount"`
	MaxResultCount int `query:"maxResultCount"`
}

const (
	DefaultMaxResultCount = 30
)

type Fields map[string]interface{}

func QueryParam(name string, ctx echo.Context) string {
	params := ctx.QueryParams()
	return params.Get(name)
}

func ReturnSuccessWithTotalCountAndItems(ctx echo.Context, status int, totalCount int64, items interface{}) error {
	return ctx.JSON(status, api.Result{
		Success: true,
		Result:  api.ArrayResult{TotalCount: totalCount, Items: items},
	})
}
func ReturnSuccessWithResult(ctx echo.Context, status int, result interface{}) error {
	return ctx.JSON(status, api.Result{
		Success: true,
		Result:  result,
	})
}

func ReturnError(ctx echo.Context, status int, apiError api.Error) error {
	return ctx.JSON(status, api.Result{
		Success: false,
		Error: api.Error{
			Code:    apiError.Code,
			Message: fmt.Sprintf("%v-%v", ProjectName, apiError.Message),
		},
	})
}

func ReturnErrorWithParameter(c echo.Context, parameters []string) error {
	return c.JSON(http.StatusBadRequest, api.Result{
		Success: false,
		Error: api.Error{
			Code:    api.ErrorParameter.Code,
			Message: fmt.Sprintf("%v-%v", ProjectName, api.ErrorParameter.Message),
			Details: fmt.Sprint(parameters),
		},
	})
}
