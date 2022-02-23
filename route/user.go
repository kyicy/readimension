package route

import (
	"fmt"
	"net/http"

	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility/config"
	"github.com/labstack/echo/v4"
)

const (
	signUpPath  = "/sign-up"
	signInPath  = "/sign-in"
	signUpFlash = "sign_up"
	signInFlash = "sign_in"
)

func getSignUp(c echo.Context) error {
	tc := newTemplateCommon(c, "Sign Up")
	tc.logout()

	s, _ := tc.GetSession()
	flashes := s.Flashes(signUpPath)

	for _, flash := range flashes {
		tc.Flashes = append(tc.Flashes, flash.(string))
	}

	s.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "user/sign_up", tc)
}

// SignUpUser binds incoming data
type signUpUser struct {
	Username  string `form:"username" validate:"required"`
	Email     string `form:"email" validate:"required,email"`
	Password  string `form:"password" validate:"required,min=5"`
	CPassword string `form:"c_password" validate:"required,min=5"`
}

func postSignUp(c echo.Context) error {
	tc := newTemplateCommon(c, "")
	s, _ := tc.GetSession()
	u := new(signUpUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	// form validation error
	err := validate.Struct(u)
	if err != nil {
		s.AddFlash(err.Error(), signUpFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signUpPath)
	}

	// email not allowed
	if !config.HasUser(u.Email) {
		s.AddFlash("Email not allowed", signUpFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signUpPath)
	}

	// email taken error
	dbUser := new(model.User)
	model.DB.Where("email = ?", u.Email).First(&dbUser)
	if dbUser.Email == u.Email {
		s.AddFlash("Email already taken", signUpFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signUpPath)
	}

	// password not matching error
	if u.Password != u.CPassword {
		s.AddFlash("password not matching", signUpFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signUpPath)
	}

	registerUser := model.User{
		Name:     u.Username,
		Email:    u.Email,
		Password: u.Password,
	}

	model.DB.Create(&registerUser)

	list := model.List{
		Name: fmt.Sprintf("/home/%s", u.Username),
		User: registerUser.ID,
	}

	model.DB.Create(&list)
	model.DB.Model(&registerUser).Association("List").Replace(&list)
	tc.login(&registerUser)

	return c.Redirect(http.StatusSeeOther, "/")

}

func getSignIn(c echo.Context) error {
	tc := newTemplateCommon(c, "Sign In")
	tc.logout()

	s, _ := tc.GetSession()
	flashes := s.Flashes(signInFlash)

	for _, flash := range flashes {
		tc.Flashes = append(tc.Flashes, flash.(string))
	}

	s.Save(c.Request(), c.Response())
	return c.Render(http.StatusOK, "user/sign_in", tc)
}

// SignUpUser binds incoming data
type signInUser struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=5"`
}

func postSignIn(c echo.Context) error {
	tc := newTemplateCommon(c, "")
	s, _ := tc.GetSession()
	u := new(signInUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	// form validation
	if err := validate.Struct(u); err != nil {
		s.AddFlash(err.Error(), signInFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signInPath)
	}

	dbUser := new(model.User)
	model.DB.Where("email = ?", u.Email).First(&dbUser)
	if dbUser.Email != u.Email {
		s.AddFlash("Account not exist", signInFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signInPath)
	}

	if !dbUser.ValidatePassword(u.Password) {
		s.AddFlash("Password error!", signInFlash)
		s.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, signInPath)
	}

	tc.login(dbUser)
	return c.Redirect(http.StatusSeeOther, "/")
}

func getSignOut(c echo.Context) error {
	tc := newTemplateCommon(c, "")
	tc.logout()
	return c.Redirect(http.StatusFound, "/")
}
