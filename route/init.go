package route

import (
	"encoding/gob"
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	validator "github.com/go-playground/validator/v10"
	mw "github.com/kyicy/readimension/middleware"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/markbates/pkger"
)

type userData map[string]string

var validate *validator.Validate

func init() {
	gob.Register(&userData{})
	validate = validator.New()
}

// Register registers all handler to a url path
func Register(e *echo.Echo) {
	render := getRender()
	e.Renderer = render

	e.GET("/", getExplorerRoot)
	e.GET("/u/explorer", getExplorerRoot)
	e.GET("/u/explorer/:list_id", getExplorer)

	e.GET("/sign-up", getSignUp)
	e.POST("/sign-up", postSignUp)

	e.GET("/sign-in", getSignIn)
	e.POST("/sign-in", postSignIn)
	e.GET("/sign-out", getSignOut)

	conf := config.Get()
	if conf.ServeStatic {
		e.Static("/covers", filepath.Join(conf.WorkDir, "covers"))
		e.Static("/books", filepath.Join(conf.WorkDir, "books"))
	}

	userGroup := e.Group("/u", mw.UserAuth)
	userGroup.DELETE("/explorer/:list_id", deleteExplorer)
	userGroup.POST("/:list_id/books/new", postBooksNew)
	userGroup.POST("/:list_id/books/new/chunksdone", postChunksDone)
	userGroup.POST("/lists/:id/child/new", postListChildNew)

	pkger.Walk("/bib", func(fullPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		tt := strings.Split(fullPath, ":/bib")
		if len(tt) < 2 {
			return nil
		}
		path := tt[1]
		extName := filepath.Ext(path)
		mt := mime.TypeByExtension(extName)

		webPath := filepath.Join("/u", path)
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
	e.GET("/u/i", func(c echo.Context) error {
		tc := newTemplateCommon(c, "")
		return c.Render(http.StatusOK, "bibi", tc)
	})

	e.GET("/u/i/", func(c echo.Context) error {
		tc := newTemplateCommon(c, "")
		return c.Render(http.StatusOK, "bibi", tc)
	})
}

func getSessionUserID(c echo.Context) (string, error) {
	s, err := session.Get("session", c)
	if err != nil {
		return "", err
	}
	if ud, flag := s.Values["userData"]; flag {
		userDataPointer := ud.(*userData)
		userID := (*userDataPointer)["id"]
		return userID, nil
	}
	return "", errors.New("session not found")

}
