package route

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kyicy/readimension/model"
	"github.com/kyicy/readimension/utility"
	"github.com/labstack/echo"
)

func getSignUp(c echo.Context) error {
	tc := newTemplateCommon(c, "Sign Up")
	tc.logout()

	sess, _ := tc.GetSession()
	flashes := sess.Flashes("sign_up")

	for _, flash := range flashes {
		tc.Flashes = append(tc.Flashes, flash.(string))
	}

	sess.Save(c.Request(), c.Response())

	return c.Render(http.StatusOK, "user/sign_up", tc)
}

// SignUpUser binds incomming data
type signUpUser struct {
	Username  string `form:"username" validate:"required"`
	Email     string `form:"email" validate:"required,email"`
	Password  string `form:"password" validate:"required,min=5"`
	CPassword string `form:"c_password" validate:"required,min=5"`
}

func postSignUp(c echo.Context) error {
	tc := newTemplateCommon(c, "")
	sess, _ := tc.GetSession()
	u := new(signUpUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	// form validation error
	err := validate.Struct(u)
	if err != nil {
		sess.AddFlash(err.Error(), "sign_up")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/sign-up")
	}

	// email taken error
	dbUser := new(model.User)
	model.DB.Where("email = ?", u.Email).First(&dbUser)
	if dbUser.Email == u.Email {
		sess.AddFlash("Email already taken", "sign_up")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/sign-up")
	}

	// password not matching error
	if u.Password != u.CPassword {
		sess.AddFlash("password not matching", "sign_up")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/sign-up")
	}

	// Start Register Process
	// Generate a uuid, and store related information to redis, expires in one day

	v4UUID, _ := uuid.NewRandom()
	utility.RedisClient.HMSet(v4UUID.String(), map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
	})
	utility.RedisClient.Expire(v4UUID.String(), time.Minute*30)

	utility.Postman.Send(u.Username, u.Email, v4UUID.String())

	return c.Redirect(http.StatusSeeOther, "/to-be-activated")
}

type _getActivateData struct {
	IsSuccess bool
	Message   string
	*TempalteCommon
}

func getActivate(c echo.Context) error {
	data := new(_getActivateData)
	tc := newTemplateCommon(c, "Activate")
	data.TempalteCommon = tc

	uuid := c.Param("uuid")
	redisData := utility.RedisClient.HGetAll(uuid)
	val := redisData.Val()
	username := val["username"]
	email := val["email"]
	password := val["password"]

	// uuid expired or wrong uuid
	if len(username) == 0 {
		data.Message = "Link expired, please sign up again"
		return c.Render(http.StatusOK, "user/activate", data)
	}

	// email taken error
	dbUser := new(model.User)
	model.DB.Where("email = ?", email).First(&dbUser)
	if dbUser.Email == email {
		data.Message = "Email already activated"
		return c.Render(http.StatusOK, "user/activate", data)
	}

	registerUser := model.User{
		Name:     username,
		Email:    email,
		Password: password,
	}

	model.DB.Create(&registerUser)

	tc.login(&registerUser)
	data.Message = fmt.Sprintf("Welcome! %s(%s). May the force be with you", username, email)
	data.IsSuccess = true

	utility.RedisClient.Del(uuid)

	return c.Render(http.StatusOK, "user/activate", data)
}

func getSignIn(c echo.Context) error {
	tc := newTemplateCommon(c, "Sign In")
	tc.logout()

	sess, _ := tc.GetSession()
	flashes := sess.Flashes("sign_in")

	for _, flash := range flashes {
		tc.Flashes = append(tc.Flashes, flash.(string))
	}

	sess.Save(c.Request(), c.Response())
	return c.Render(http.StatusOK, "user/sign_in", tc)
}

// SignUpUser binds incomming data
type signInUser struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=5"`
}

func postSignIn(c echo.Context) error {
	tc := newTemplateCommon(c, "")
	sess, _ := tc.GetSession()
	u := new(signInUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	// form validation
	if err := validate.Struct(u); err != nil {
		sess.AddFlash(err.Error(), "sign_in")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/sign-in")
	}

	dbUser := new(model.User)
	model.DB.Where("email = ?", u.Email).First(&dbUser)
	if dbUser.Email != u.Email {
		sess.AddFlash("Account not exist", "sign_in")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/sign-in")
	}

	if !dbUser.ValidatePassword(u.Password) {
		sess.AddFlash("Password error!", "sign_in")
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusSeeOther, "/sign-in")
	}

	tc.login(dbUser)
	return c.Redirect(http.StatusSeeOther, "/")
}

func getSignOut(c echo.Context) error {
	tc := newTemplateCommon(c, "")
	tc.logout()
	return c.Redirect(http.StatusSeeOther, "/")
}

func getToBeActivated(c echo.Context) error {
	tc := newTemplateCommon(c, "To Be Activated")
	return c.Render(http.StatusOK, "user/to_be_activated", tc)
}
