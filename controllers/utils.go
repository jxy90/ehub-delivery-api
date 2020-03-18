package controllers

import (
	"fmt"
	"net/http"

	"github.com/hublabs/common/api"

	"github.com/labstack/echo"
)

var ProjectName string = "[ehub-delivery-api]"

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

func ReturnApiSucc(ctx echo.Context, status int, totalCount int64, items interface{}) error {
	return ctx.JSON(status, api.Result{
		Success: true,
		Result:  api.ArrayResult{TotalCount: totalCount, Items: items},
	})
}
func ReturnResultApiSucc(ctx echo.Context, status int, result interface{}) error {
	return ctx.JSON(status, api.Result{
		Success: true,
		Result:  result,
	})
}

func ReturnApiWarn(ctx echo.Context, status int, apiError api.Error, err error) error {
	str := ""
	if err != nil {
		str = fmt.Sprint(err)
	}

	return ctx.JSON(status, api.Result{
		Success: false,
		Error: api.Error{
			Code:    apiError.Code,
			Message: fmt.Sprintf(apiError.Message),
			Details: ProjectName + str,
		},
	})
}

func ReturnApiParameterWarn(c echo.Context, parameters []string) error {
	return c.JSON(http.StatusBadRequest, api.Result{
		Success: false,
		Error: api.Error{
			Code:    api.ErrorParameter.Code,
			Message: fmt.Sprintf(api.ErrorParameter.Message),
			Details: ProjectName + fmt.Sprint(parameters),
		},
	})
}
