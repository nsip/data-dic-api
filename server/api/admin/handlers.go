package admin

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	. "github.com/digisan/go-generics/v2"
	gm "github.com/digisan/go-mail"
	structtool "github.com/digisan/gotk/struct-tool"
	lk "github.com/digisan/logkit"
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/db"
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
// @Param   fields  path  string false "which user's fields (sep by ',') want to list. if empty, return all fields"
// @Success 200 "OK - list successfully"
// @Failure 401 "Fail - unauthorized error"
// @Failure 403 "Fail - forbidden error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/user/list/{fields} [get]
// @Security ApiKeyAuth
func ListUser(c echo.Context) error {

	lk.Log("Enter: ListUser")

	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		admin   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(admin)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusForbidden, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
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
		fields = c.Param("fields")
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

	if len(fields) == 0 {
		return c.JSON(http.StatusOK, users)
	}

	if fields, err = url.QueryUnescape(fields); err != nil {
		c.String(http.StatusBadRequest, "'fields' is invalid")
	}

	fieldsUser := []string{}
	for _, field := range strings.Split(fields, ",") {
		fieldsUser = AppendIf(In(field, "uname", "Uname", "ID", "Id", "id"), fieldsUser, "UName")
		fieldsUser = AppendIf(In(field, "email", "Email"), fieldsUser, "Email")
		fieldsUser = AppendIf(In(field, "name", "Name"), fieldsUser, "Name")
	}
	rt := FilterMap(users, nil, func(i int, e *u.User) any {
		v, err := structtool.PartialAsMap(e, fieldsUser...)
		lk.WarnOnErr("%v", err)
		return v
	})
	return c.JSON(http.StatusOK, rt)
}

// @Title list user's action record
// @Summary list user's action record
// @Description
// @Tags    Admin
// @Accept  json
// @Produce json
// @Param   uname  query string true "user registered unique name"
// @Param   action path  string true "which action type [submit, approve, subscribe] record want to list"
// @Success 200 "OK - list successfully"
// @Failure 401 "Fail - unauthorized error"
// @Failure 403 "Fail - forbidden error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/user/action-list/{action} [get]
// @Security ApiKeyAuth
func ListUserAction(c echo.Context) error {

	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		admin   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(admin)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusForbidden, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	// if user.MemLevel != 3 {
	// 	return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	// }

	// --- //

	var (
		uname  = c.QueryParam("uname") // other uname
		action = c.Param("action")     // action type: submit, approve, subscribe
	)

	ls, err := db.ListActionRecord(uname, db.DbColType(action))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ls)
}

// @Title
// @Summary
// @Description
// @Tags    Admin
// @Accept  multipart/form-data
// @Produce json
// @Param   unames  formData string true "unique user names, separator is ',' "
// @Param   subject formData string true "subject for email"
// @Param	body    formData string true "body for email"
// @Success 200 "OK - list successfully"
// @Failure 401 "Fail - unauthorized error"
// @Failure 403 "Fail - forbidden error"
// @Failure 500 "Fail - internal error"
// @Router /api/admin/email [post]
// @Security ApiKeyAuth
func SendEmail(c echo.Context) error {

	var (
		userTkn = c.Get("user").(*jwt.Token)
		claims  = userTkn.Claims.(*u.UserClaims)
		admin   = claims.UName
	)

	user, ok, err := u.LoadActiveUser(admin)

	switch {
	case err != nil:
		return c.String(http.StatusInternalServerError, err.Error())
	case !ok:
		return c.String(http.StatusForbidden, fmt.Sprintf("invalid user status@[%s], dormant?", user.UName))
	}

	// if user.MemLevel != 3 {
	// 	return c.String(http.StatusUnauthorized, "failed, you are not authorized to this api")
	// }

	const (
		sep = "," // separator for unames
	)

	var (
		unames  = c.FormValue("unames")  // recipients, separator is ','
		subject = c.FormValue("subject") // email title
		body    = c.FormValue("body")    // email content
	)

	type retType struct {
		OK     bool
		Sent   []string
		Failed []string
		Err    []error
	}
	ret := []retType{}

	for _, uname := range strings.Split(unames, sep) {
		lk.Log("[%v] [%v] [%v]", uname, subject, body)

		user, ok, err = u.LoadUser(uname, true)
		switch {
		case err != nil:
			return c.String(http.StatusInternalServerError, err.Error())
		case !ok:
			return c.String(http.StatusBadRequest, fmt.Sprintf("[%s] doesn't exist", uname))
		}

		ok, sent, failed, errs := gm.SendMail(subject, body, user.Email)
		ret = append(ret, retType{
			OK:     ok,
			Sent:   sent,
			Failed: failed,
			Err:    errs,
		})
	}

	return c.JSON(http.StatusOK, ret)
}
