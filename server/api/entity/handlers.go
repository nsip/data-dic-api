package entity

import (
	"net/http"

	. "github.com/digisan/go-generics/v2"
	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/db"
)

// @Title insert or update one entity data
// @Summary insert or update one entity data by a json file
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Param   dbName  query string false "database name"    default(dictionary)
// @Param   colName query string false "collection name"  default(entity)
// @Param   entity  path  string true  "entity name for incoming entity data"
// @Param   data    body  string true  "entity json data for uploading" Format(binary)
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/insert [post]
func Insert(c echo.Context) error {

	var (
		dbName  = c.QueryParam("db")
		colName = c.QueryParam("col")
		entity  = c.Param("entity")
		dataRdr = c.Request().Body
	)

	dbName = IF(len(dbName) == 0, "dictionary", dbName)
	colName = IF(len(colName) == 0, "entity", colName)

	if len(entity) == 0 {
		return c.String(http.StatusBadRequest, "entity name is empty")
	}
	if dataRdr == nil {
		return c.String(http.StatusBadRequest, "payload for insert is empty")
	}

	db.UseDbCol(dbName, colName)

	id, err := db.Insert(dataRdr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, id)
}

// @Title find entity json
// @Summary find entities json content by pass a json query string via payload
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Param   dbName  query string false "database name"
// @Param   colName query string false "collection name"
// @Param   data    body  string true  "json data for query" Format(binary)
// @Success 200 "OK - find successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/find [get]
func Find(c echo.Context) error {

	var (
		dbName  = c.QueryParam("db")
		colName = c.QueryParam("col")
		qryRdr  = c.Request().Body
	)

	dbName = IF(len(dbName) == 0, "dictionary", dbName)
	colName = IF(len(colName) == 0, "entity", colName)

	if qryRdr == nil {
		return c.String(http.StatusBadRequest, "payload for query is empty")
	}

	db.UseDbCol(dbName, colName)

	results, err := db.Find[EntityType](qryRdr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}
