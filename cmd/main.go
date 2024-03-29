package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/PODO/datepath-server"
	_ "github.com/PODO/datepath-server/docs"
	"github.com/PODO/datepath-server/handler/direction"
	"github.com/PODO/datepath-server/handler/local"
	"github.com/PODO/datepath-server/pkg/logrotater"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Swagger Date Path API
// @version 1.0
// @description 데이트 장소 만들기 위한 API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
func main() {
	confFilePath := flag.String("conf", "", "config file path")
	echoLogPath := flag.String("log-path", "", "echo output log path")
	echoLogLevel := flag.String("level", "info", "echo log level")
	flag.Parse()

	e := echo.New()
	e.HideBanner = true
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
	)

	conf, err := datepath_server.LoadConfig(*confFilePath)
	if err != nil {
		panic(err)
	}
	e.Use(bindConfig(conf))

	clients := datepath_server.NewServiceClients(conf)
	e.Use(bindClients(clients))

	logFileName := "echo.log"
	logFilePath := fmt.Sprintf("%s%s", *echoLogPath, logFileName)
	if *echoLogLevel == "debug" {
		e.Logger.SetLevel(log.DEBUG)
	} else {
		e.Logger.SetLevel(log.INFO)
	}
	lr := logrotater.New(logFilePath, &conf.LogRotater)
	e.Logger.SetOutput(&lr.Logger)
	go lr.Run()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	lg := e.Group("/local/search/keyword")
	lg.POST("", local.SearchKeywordHandler)

	dg := e.Group("/direction/coordinates")
	dg.GET("", direction.DirectionHandler)

	servicePort := fmt.Sprintf(":%d", conf.Server.Port)
	go func() {
		if err := e.Start(servicePort); err != nil {
			if err == http.ErrServerClosed {
				e.Logger.Infof("Echo Start: %v", err)
			} else {
				e.Logger.Fatalf("Echo Start: %v", err)
			}
		}
	}()
	if conf.Server.HttpsOn {
		go func() {
			if err := e.StartTLS(":443", conf.Server.SSLCrtPath, conf.Server.SSLKeyPath); err != nil {
				if err == http.ErrServerClosed {
					e.Logger.Infof("Echo Start: %v", err)
				} else {
					e.Logger.Fatalf("Echo Start: %v", err)
				}
			}
		}()
	}

	// 종료 시그널 대기
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-shutdown

	timeout := time.Duration(conf.Server.ShutdownTimeoutInSecond) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	e.Logger.Infof("Shutdown start timeout(%s)", timeout)
	e.Logger.Infof("Shutdown server start")
	err = e.Shutdown(ctx)
	e.Logger.Infof("Shutdown server finished: %v", err)
}

func bindConfig(conf *datepath_server.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("datepathServerConfig", conf)
			return next(c)
		}
	}
}

func bindClients(clients *datepath_server.ServiceClients) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("datepathServiceClients", clients)
			return next(c)
		}
	}
}
