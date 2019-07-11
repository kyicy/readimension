package middleware

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// UserAuth bla bla
func UserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)

		if sess.Values["userExist?"] == true {
			return next(c)
		}
		return c.Redirect(http.StatusSeeOther, "/sign-in")
	}
}
