package route

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type TempalteCommon struct {
	echo.Context
	Title           string
	Active          string
	Flashes         []string
	GoogleAnalytics string
	GoogleAdsense   string
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

// HasGA aka check if user has configured google analytics
func (tc *TempalteCommon) HasGoogleAnlytics() bool {
	return len(tc.GoogleAnalytics) > 0
}

func (tc *TempalteCommon) HasGoogleAdsense() bool {
	return len(tc.GoogleAdsense) > 0
}

func newTemplateCommon(c echo.Context, title string) *TempalteCommon {
	title = title + " - Readimension"
	configuration := config.Get()
	return &TempalteCommon{
		Context:         c,
		Title:           title,
		Active:          c.Request().URL.Path,
		GoogleAnalytics: configuration.GoogleAnalytics,
		GoogleAdsense:   configuration.GoogleAdsense,
	}
}
