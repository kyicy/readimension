package route

import (
	"encoding/gob"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	mw "github.com/kyicy/readimension/middleware"
	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"gopkg.in/go-playground/validator.v9"
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

	e.GET("/sign-up", getSignUp)
	e.POST("/sign-up", postSignUp)

	e.GET("/sign-in", getSignIn)
	e.POST("/sign-in", postSignIn)
	e.GET("/sign-out", getSignOut)

	e.GET("/", getExplorerRoot, mw.UserAuth)

	conf := config.Get()
	if conf.ServeStatic {
		e.Static("/covers", "covers")
		e.Static("/books", "books")
	}

	userGroup := e.Group("/u", mw.UserAuth)
	userGroup.GET("/explorer", getExplorerRoot)
	userGroup.DELETE("/explorer/:list_id", deleteExplorer)
	userGroup.GET("/explorer/:list_id", getExplorer)
	userGroup.GET("/explorer", getExplorerRoot)
	userGroup.POST("/:list_id/books/new", postBooksNew)
	userGroup.POST("/:list_id/books/new/chunksdone", postChunksDone)
	userGroup.POST("/lists/:id/child/new", postListChildNew)

	box := packr.NewBox("../bib")
	box.Walk(func(path string, f packr.File) error {
		extName := filepath.Ext(path)
		mt := mime.TypeByExtension(extName)

		userGroup.GET("/"+path, func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=3600")
			r := strings.NewReader(box.String(path))
			return c.Stream(http.StatusOK, mt, r)
		})
		return nil
	})
	userGroup.GET("/i", func(c echo.Context) error {
		tc := newTemplateCommon(c, "")
		return c.Render(http.StatusOK, "bibi", tc)
	})

	userGroup.GET("/i/", func(c echo.Context) error {
		tc := newTemplateCommon(c, "")
		return c.Render(http.StatusOK, "bibi", tc)
	})
}

func getSessionUser(c echo.Context) (*model.User, error) {
	userID, err := getSessionUserID(c)
	if err != nil {
		return nil, err
	}
	var userRecord model.User
	model.DB.Where("id = ?", userID).First(&userRecord)
	return &userRecord, nil
}

func getSessionUserID(c echo.Context) (string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return "", err
	}
	ud := sess.Values["userData"]
	userDataPointer := ud.(*userData)
	userID := (*userDataPointer)["id"]
	return userID, nil
}
