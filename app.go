package main

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2"
	"github.com/kyicy/readimension/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func createInstance(env string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
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
	} else {
		e.Static("/", "public")
	}

	route.Register(e)

	return e
}

func bundle(e *echo.Echo) {
	box := packr.New("public_box", "./public")

	box.Walk(func(path string, f packr.File) error {

		extName := filepath.Ext(path)
		mt := mime.TypeByExtension(extName)

		e.GET("/"+path, func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=3600")
			s, _ := box.FindString(path)
			r := strings.NewReader(s)
			return c.Stream(http.StatusOK, mt, r)
		})
		return nil
	})
}
