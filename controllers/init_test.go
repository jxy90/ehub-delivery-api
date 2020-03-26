package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/hublabs/common/auth"
	"github.com/hublabs/delivery-api/models"

	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
	dbContext        context.Context
)

func TestMain(m *testing.M) {
	db := enterTest()
	code := m.Run()
	exitTest(db)
	os.Exit(code)
}

func enterTest() *xorm.Engine {
	runtime.GOMAXPROCS(1)
	//xormEngine, err := xorm.NewEngine("sqlite3", ":memory:"); if err != nil {
	//		panic(err)
	//	}
	xormEngine, err := xorm.NewEngine("mysql", "root:@tcp(localhost:3306)/hublabs_delivery?charset=utf8")
	if err != nil {
		panic(err)
	}
	dbContext = context.WithValue(context.Background(), echomiddleware.ContextDBName, xormEngine.NewSession())
	if err != nil {
		panic(err)
	}
	if err = models.DropTables(xormEngine); err != nil {
		panic(err)
	}
	if err = models.Init(xormEngine); err != nil {
		panic(err)
	}
	if err := CreateSeedData(xormEngine); err != nil {
		fmt.Println("create seed data err:", err)
	}

	echoApp = echo.New()
	echoApp.Validator = &Validator{}

	db := echomiddleware.ContextDB("test", xormEngine, echomiddleware.KafkaConfig{})
	jwt := middleware.JWT([]byte(os.Getenv("JWT_SECRET")))
	behaviorLogger := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				behaviorlog.LogContextName, behaviorlog.New("test", req),
			)))
			return next(c)
		}
	}
	jwtutil.SetJwtSecret(os.Getenv("JWT_SECRET"))
	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return behaviorLogger(jwt(auth.UserClaimMiddleware()(db(handlerFunc))))(c)
	}
	return xormEngine
}

func exitTest(db *xorm.Engine) {
	// if err := models.DropTables(db); err != nil {
	// 	panic(err)
	// }
}

func setHeader(r *http.Request) {
	token, _ := jwtutil.NewToken(map[string]interface{}{"aud": "colleague", "tenantCode": "test", "iss": "colleague"})
	r.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
