package main

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	mw "github.com/kyicy/readimension/middleware"
	"github.com/kyicy/readimension/route"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func createInstance(env string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         3600,
	}))

	isProduction := env == "production"
	if isProduction {
		bundle(e)
		e.Use(mw.Minify)
	} else {
		e.Static("/", "public")
		e.Static("/covers", "covers")
		e.Static("/books", "books")
	}

	route.Register(e)

	return e
}

func bundle(e *echo.Echo) {
	box := packr.NewBox("./public")

	box.Walk(func(path string, f packr.File) error {

		extName := filepath.Ext(path)
		mt := mime.TypeByExtension(extName)

		e.GET("/"+path, func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=3600")
			r := strings.NewReader(box.String(path))
			return c.Stream(http.StatusOK, mt, r)
		})
		return nil
	})
}
