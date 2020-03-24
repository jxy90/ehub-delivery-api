package main

import (
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/hublabs/ehub-delivery-api/config"
	"github.com/hublabs/ehub-delivery-api/controllers"
	"github.com/hublabs/ehub-delivery-api/models"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/urfave/cli/v2"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	c := config.Init(os.Getenv("APP_ENV"))
	xormEngine, err := xorm.NewEngine(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}

	defer xormEngine.Close()

	app := cli.NewApp()
	app.Name = "ehub-delivery-api"
	app.Commands = []*cli.Command{
		{
			Name:  "api-server",
			Usage: "run api server",
			Action: func(cli *cli.Context) error {
				if err := initEchoApp(xormEngine).Start(":" + c.HttpPort); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:  "seed",
			Usage: "create seed data",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func initEchoApp(xormEngine *xorm.Engine) *echo.Echo {
	xormEngine.SetMaxOpenConns(50)
	xormEngine.SetMaxIdleConns(50)
	xormEngine.SetConnMaxLifetime(60 * time.Second)

	e := echo.New()
	r := echoswagger.New(e, "doc", &echoswagger.Info{
		Title:       "EHub Delivery API",
		Description: "This is docs for ehub-delivery-api",
		Version:     "1.0.0",
	})
	InitControllers(e)

	e.Static("/static", "static")
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())

	c := config.Init(os.Getenv("APP_ENV"))
	db := initDB(c.Database.Driver, c.Database.Connection)

	if c.Env != "production" {
		if err := models.Init(db); err != nil {
			panic(err)
		}
	}
	e.Use(echomiddleware.ContextDB(c.ServiceName, db, c.Database.Logger.Kafka))

	controllers.DeliveryController{}.Init(r.Group("Delivery", "/v1/delivery"))
	return e
}

func initDB(driver, connection string) *xorm.Engine {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		panic(err)
	}

	if driver == "sqlite3" {
		runtime.GOMAXPROCS(1)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Minute * 10)

	env := os.Getenv("APP_ENV")
	if env == "staging" || env == "" {
		db.ShowSQL()
	}
	return db
}

func InitControllers(e *echo.Echo) {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
}
