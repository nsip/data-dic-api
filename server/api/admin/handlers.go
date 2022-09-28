package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// @Title list users' info
// @Summary list users' info
// @Description
// @Tags    Admin
// @Accept  json
// @Produce json
// @Param   uname  query string false "user filter with uname wildcard(*)"
// @Param   name   query string false "user filter with name wildcard(*)"
// @Param   active query string false "user filter with active status"
// @Param   rType  path  string false "which user's field want to list"
// @Success 200 "OK - list successfully"
// @Failure 401 "Fail - unauthorized error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/user/list/{rType} [get]
// @Security ApiKeyAuth
func ListUser(c echo.Context) error {
	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		uname   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(uname)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusInternalServerError, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	// if user.MemLevel != 3 {
	// 	return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	// }

	// --- //

	var (
		active = c.QueryParam("active")
		wUname = c.QueryParam("uname")
		wName  = c.QueryParam("name")
		rUname = wc2re(wUname)
		rName  = wc2re(wName)
		rtType = c.Param("rType")
	)

	users, err := u.ListUser(func(u *u.User) bool {
		switch {
		case len(wUname) > 0 && !rUname.MatchString(u.UName):
			return false
		case len(wName) > 0 && !rName.MatchString(u.Name):
			return false
		case len(active) > 0:
			if bActive, err := strconv.ParseBool(active); err == nil {
				return bActive == u.Active
			}
			return false
		default:
			return true
		}
	})

	for _, user := range users {
		user.Password = strings.Repeat("*", len(user.Password))
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var rt any
	switch rtType {
	case "uname", "Uname", "ID", "Id", "id":
		rt = FilterMap(users, nil, func(i int, e *u.User) string { return e.UName })
	case "email", "Email":
		rt = FilterMap(users, nil, func(i int, e *u.User) string { return e.Email })
	case "name", "Name":
		rt = FilterMap(users, nil, func(i int, e *u.User) string { return e.Name })
	default:
		rt = users
	}
	return c.JSON(http.StatusOK, rt)
}
