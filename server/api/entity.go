package api

import (
	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/entity"
)

// register to main echo Group

// "/api/entity"
func EntityHandler(r *echo.Group) {

	var mGET = map[string]echo.HandlerFunc{
		// "/find": entity.Find,
		"/entities":  entity.AllEntities,
		"/list_names": entity.AllEntityNames,
	}
	var mPOST = map[string]echo.HandlerFunc{
		"/upsert/:valType": entity.Upsert,
	}
	var mPUT = map[string]echo.HandlerFunc{}
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
		"GET":    r.GET,
		"POST":   r.POST,
		"PUT":    r.PUT,
		"DELETE": r.DELETE,
		"PATCH":  r.PATCH,
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
