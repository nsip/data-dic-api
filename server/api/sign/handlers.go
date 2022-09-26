package sign

import (
	"net/http"
	"sync"
	"time"

	lk "github.com/digisan/logkit"
	si "github.com/digisan/user-mgr/sign-in"
	su "github.com/digisan/user-mgr/sign-up"
	u "github.com/digisan/user-mgr/user"
	"github.com/labstack/echo/v4"
)

// *** after implementing, register with path in 'sign.go' *** //

var (
	MapUserClaims = &sync.Map{} // map[string]*u.UserClaims, *** record logged-in user claims  ***
)

// @Title register a new user
// @Summary sign up action, send user's basic info for registry
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname   formData   string  true  "unique user name"
// @Param   email   formData   string  true  "user's email" Format(email)
// @Param   pwd     formData   string  true  "user's password"
// @Success 200 "OK - then waiting for verification code"
// @Failure 400 "Fail - invalid registry fields"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/new [post]
func NewUser(c echo.Context) error {

	// lk.Debug("[%v] [%v] [%v]", c.FormValue("uname"), c.FormValue("email"), c.FormValue("pwd"))

	user := &u.User{
		Core: u.Core{
			UName:    c.FormValue("uname"),
			Email:    c.FormValue("email"),
			Password: c.FormValue("pwd"),
		},
		Profile: u.Profile{
			Name:           "",
			Phone:          "",
			Country:        "",
			City:           "",
			Addr:           "",
			PersonalIDType: "",
			PersonalID:     "",
			Gender:         "",
			DOB:            "",
			Position:       "",
			Title:          "",
			Employer:       "",
			Bio:            "",
			AvatarType:     "",
			Avatar:         []byte{},
		},
		Admin: u.Admin{
			Regtime:   time.Now().Truncate(time.Second),
			Active:    true,
			Certified: false,
			Official:  false,
			SysRole:   "",
			MemLevel:  0,
			MemExpire: time.Time{},
			Tags:      "",
		},
	}

	// su.SetValidator(map[string]func(string) bool{ })

	lk.Log("%v", user)

	if err := su.ChkInput(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	///////////////////////////////////////////////
	// simple sing up, ignore email verification //
	///////////////////////////////////////////////

	// store into db
	if err := su.Store(user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// sign-up ok calling...
	{
	}

	return c.JSON(http.StatusOK, "registered successfully")
}

// @Title sign in
// @Summary sign in action. if ok, got token
// @Description
// @Tags    Sign
// @Accept  multipart/form-data
// @Produce json
// @Param   uname formData string true "user name or email"
// @Param   pwd   formData string true "password" Format(password)
// @Success 200 "OK - sign-in successfully"
// @Failure 400 "Fail - incorrect password"
// @Failure 500 "Fail - internal error"
// @Router /api/sign/in [post]
func LogIn(c echo.Context) error {

	var (
		uname = c.FormValue("uname")
		pwd   = c.FormValue("pwd")
		email = c.FormValue("uname")
	)

	lk.Debug("login: [%v] [%v]", uname, pwd)

	user := &u.User{
		Core: u.Core{
			UName:    uname,
			Password: pwd,
			Email:    email,
		},
		Profile: u.Profile{},
		Admin:   u.Admin{},
	}

	if err := si.CheckUserExisting(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if !si.PwdOK(user) { // if successful, user updated.
		return c.String(http.StatusBadRequest, "incorrect password")
	}

	// fmt.Println(user)

	// now, user is real user in db
	defer lk.FailOnErr("%v", si.Trail(user.UName)) // Refresh Online Users, here UName is real

	// log in ok calling...
	{
	}

	claims := u.MakeUserClaims(user)
	defer func() { MapUserClaims.Store(user.UName, claims) }() // save current user claims for other usage

	token := claims.GenToken()
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
		"auth":  "Bearer " + token,
	})
}
