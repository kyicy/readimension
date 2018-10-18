package route

import (
	"encoding/gob"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/sessions"
	mw "github.com/kyicy/readimension/middleware"
	"github.com/kyicy/readimension/model"
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

type TempalteCommon struct {
	echo.Context
	Title   string
	Active  string
	Flashes []string
}

func (tc *TempalteCommon) GetSession() (*sessions.Session, error) {
	return session.Get("session", tc.Context)
}

func (tc *TempalteCommon) logout() {
	sess, _ := tc.GetSession()
	sess.Values["userExist?"] = false
	delete(sess.Values, "userData")
	sess.Save(tc.Request(), tc.Response())
}

func (tc *TempalteCommon) login(u *model.User) {
	sess, _ := tc.GetSession()
	sess.Values["userExist?"] = true
	sess.Values["userData"] = userData{
		"id":    fmt.Sprintf("%d", u.ID),
		"name":  u.Name,
		"email": u.Email,
	}
	sess.Save(tc.Request(), tc.Response())
}

func newTemplateCommon(c echo.Context, title string) *TempalteCommon {
	title = title + " - Readimension"
	return &TempalteCommon{
		Context: c,
		Title:   title,
		Active:  c.Request().URL.Path,
	}
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

		e.GET("/u/"+path, func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "max-age=3600")
			r := strings.NewReader(box.String(path))
			return c.Stream(http.StatusOK, mt, r)
		})
		return nil
	})
	e.GET("/u/i", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=3600")
		r := box.String("i/index.html")
		return c.HTML(http.StatusOK, r)
	})

	e.GET("/u/i/", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=3600")
		r := box.String("i/index.html")
		return c.HTML(http.StatusOK, r)
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
