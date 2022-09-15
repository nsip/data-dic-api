package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nsip/data-dic-api/server/api"
	_ "github.com/nsip/data-dic-api/server/docs" // once `swag init`, comment it out
	"github.com/postfinance/single"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var (
	fHttp2 = false
)

func init() {
	lk.WarnDetail(false)
}

// @title National Education Data Dictionary API
// @version 1.0
// @description This is national education data dictionary backend-api server. Updated@ 2022-09-15T09:29:03+10:00
// @termsOfService
// @contact.name API Support
// @contact.url
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:1323
// @BasePath
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name authorization
func main() {

	http2Ptr := flag.Bool("http2", false, "http2 mode?")
	flag.Parse()
	fHttp2 = *http2Ptr

	// only one instance
	const dir = "./tmp-locker"
	gio.MustCreateDir(dir)
	one, err := single.New("echo-service", single.WithLockPath(dir))
	lk.FailOnErr("%v", err)
	lk.FailOnErr("%v", one.Lock())
	defer func() {
		lk.FailOnErr("%v", one.Unlock())
		os.RemoveAll(dir)
		lk.Log("Server Exited Successfully")
	}()

	// start Service
	done := make(chan string)
	echoHost(done)
	lk.Log(<-done)
}

func waitShutdown(e *echo.Echo) {
	go func() {
		// defer Close Database // after closing echo, close db

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		lk.Log("Got Ctrl+C")

		// other clean-up before closing echo
		{

		}

		// shutdown echo
		lk.FailOnErr("%v", e.Shutdown(ctx)) // close echo at e.Shutdown
	}()
}

func echoHost(done chan<- string) {
	go func() {
		defer func() { done <- "Echo Shutdown Successfully" }()

		e := echo.New()
		defer e.Close()

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.BodyLimit("2G"))
		// CORS
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
			AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}))

		// waiting for shutdown
		waitShutdown(e)

		// host swagger http://localhost:1323/swagger/index.html
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		// groups without middleware
		{
			// api.SignHandler(e.Group("/api/xxx"))
		}

		// other groups with middleware
		groups := []string{
			"/api/dictionary",
		}
		handlers := []func(*echo.Group){
			api.Handler,
		}
		for i, group := range groups {
			r := e.Group(group)
			r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					return next(c)
				}
			})
			handlers[i](r)
		}

		// running...
		var err error
		if fHttp2 {
			err = e.StartTLS(":1323", "./cert/public.pem", "./cert/private.pem")
		} else {
			err = e.Start(":1323")
		}
		lk.FailOnErrWhen(err != http.ErrServerClosed, "%v", err)
	}()
}
