package entity

import (
	"net/http"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	"github.com/labstack/echo/v4"
	lk "github.com/digisan/logkit"
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
// @Router /api/entity/insert/{entity} [post]
func Insert(c echo.Context) error {

	var (
		dbName  = c.QueryParam("dbName")
		colName = c.QueryParam("colName")
		entity  = c.Param("entity")
		dataRdr = c.Request().Body
	)

	dbName = IF(len(dbName) == 0, dbDefault, dbName)
	colName = IF(len(colName) == 0, colDefault, colName)
	if len(entity) == 0 {
		return c.String(http.StatusBadRequest, "entity name is empty")
	}
	if dataRdr == nil {
		return c.String(http.StatusBadRequest, "payload for insert is empty")
	}

	lk.Log("using database: [%s]; using collection: [%s]", dbName, colName)
	mh.UseDbCol(dbName, colName)

	id, err := mh.Insert(dataRdr)
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
// @Param   dbName  query string false "database name"       default(dictionary)
// @Param   colName query string false "collection name"     default(entity)
// @Param   data    body  string true  "json data for query" Format(binary)
// @Success 200 "OK - find successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/find [get]
func Find(c echo.Context) error {

	var (
		dbName  = c.QueryParam("dbName")
		colName = c.QueryParam("colName")
		qryRdr  = c.Request().Body
	)

	dbName = IF(len(dbName) == 0, dbDefault, dbName)
	colName = IF(len(colName) == 0, colDefault, colName)
	// if qryRdr == nil {
	// 	return c.String(http.StatusBadRequest, "payload for query is empty")
	// }

	mh.UseDbCol(dbName, colName)

	results, err := mh.Find[EntityType](qryRdr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// func AllEntities(c echo.Context) error {

// }

// func AllEntityNames(c echo.Context) error {

// }
