package route

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
	"github.com/kyicy/readimension/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"gopkg.in/go-playground/validator.v9"

	mw "github.com/kyicy/readimension/middleware"
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
		"id":    string(u.ID),
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
	e.GET("/to-be-activated", getToBeActivated)
	e.GET("/activate/:uuid", getActivate)
	e.GET("/sign-out", getSignOut)

	e.GET("/", getBooks, mw.UserAuth)

	userGroup := e.Group("/u", mw.UserAuth)
	userGroup.GET("/stream", getBooks)
	userGroup.GET("/books", getBooks)

	userGroup.GET("/books/new", getBooksNew)
	userGroup.POST("/books/new", postBooksNew)
	userGroup.POST("/books/new/chunksdone", postChunksDone)

	userGroup.GET("/lists", getLists)
	userGroup.GET("/lists/new", getListsNew)
}
