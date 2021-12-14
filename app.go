package main

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kyicy/readimension/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/pkger"
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

	bundle(e)
	route.Register(e)

	return e
}

func bundle(e *echo.Echo) {
	pkger.Walk("/public", func(fullPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		tt := strings.Split(fullPath, ":/public")
		if len(tt) < 2 {
			return nil
		}
		path := tt[1]
		extName := filepath.Ext(path)
		mt := mime.TypeByExtension(extName)

		webPath := filepath.Join("/", path)
		e.GET(webPath, func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=3600")
			f, err := pkger.Open(fullPath)
			if err != nil {
				return err
			}
			defer f.Close()
			return c.Stream(http.StatusOK, mt, f)
		})
		return nil
	})
}
