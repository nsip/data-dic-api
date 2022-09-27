package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/user"
)

// register to main echo Group

// /api/user/
func SignHandler(e *echo.Group) {

	var mGET = map[string]echo.HandlerFunc{}

	var mPOST = map[string]echo.HandlerFunc{
		"/sign-up": user.NewUser,
		"/sign-in": user.LogIn,
	}

	var mPUT = map[string]echo.HandlerFunc{
		"/sign-out": user.SignOut,
	}

	var mDELETE = map[string]echo.HandlerFunc{}
	var mPATCH = map[string]echo.HandlerFunc{}

	// ------------------------------------------------------- //

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	mRegAPIs := map[string]map[string]echo.HandlerFunc{
		"GET":    mGET,
		"POST":   mPOST,
		"PUT":    mPUT,
		"DELETE": mDELETE,
		"PATCH":  mPATCH,
		// others...
	}

	mRegMethod := map[string]func(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route{
		"GET":    e.GET,
		"POST":   e.POST,
		"PUT":    e.PUT,
		"DELETE": e.DELETE,
		"PATCH":  e.PATCH,
		// others...
	}

	for _, m := range methods {
		mAPI, method := mRegAPIs[m], mRegMethod[m]
		for path, handler := range mAPI {
			if handler == nil {
				continue
			}
			method(path, handler)
		}
	}
}
