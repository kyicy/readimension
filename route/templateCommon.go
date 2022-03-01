package route

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// TemplateCommon : info container for rendering templates
type TemplateCommon struct {
	echo.Context
	Title           string
	Active          string
	Flashes         []string
	GoogleAnalytics string
}

// GetSession : retrieve session object
func (tc *TemplateCommon) GetSession() (*sessions.Session, error) {
	return session.Get("session", tc.Context)
}

func (tc *TemplateCommon) logout() {
	s, _ := tc.GetSession()
	s.Values["userExist?"] = false
	delete(s.Values, "userData")
	s.Save(tc.Request(), tc.Response())
}

func (tc *TemplateCommon) login(u *model.User) {
	s, _ := tc.GetSession()
	s.Values["userExist?"] = true
	s.Values["userData"] = userData{
		"id":    fmt.Sprintf("%d", u.ID),
		"name":  u.Name,
		"email": u.Email,
	}
	s.Save(tc.Request(), tc.Response())
}

// HasGoogleAnalytics : check if user has configured google analytics
func (tc *TemplateCommon) HasGoogleAnalytics() bool {
	return len(tc.GoogleAnalytics) > 0
}

func newTemplateCommon(c echo.Context, title string) *TemplateCommon {
	title = title + " - Readimension"
	configuration := config.Get()
	return &TemplateCommon{
		Context:         c,
		Title:           title,
		Active:          c.Request().URL.Path,
		GoogleAnalytics: configuration.GoogleAnalytics,
	}
}
