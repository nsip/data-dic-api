package entity

import (
	"fmt"
	"net/http"

	mh "github.com/digisan/db-helper/mongo"
	// . "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
)

// @Title mongodb config for entity storage
// @Summary set mongodb database and collection for entity storage
// @Description
// @Tags    Entity
// @Accept  multipart/form-data
// @Produce json
// @Param   database   formData string true "database name"   default(dictionary)
// @Param   collection formData string true "collection name" default(entity)
// @Success 200 "OK - set db successfully"
// @Failure 400 "Fail - invalid fields"
// @Router /api/entity/db [put]
func UseDbCol(c echo.Context) error {
	err := c.Bind(&cfg)
	switch {
	case err != nil:
		return c.String(http.StatusBadRequest, fmt.Sprintf("bad request for binding DbConfig: %v", err))
	case len(cfg.Database) == 0:
		return c.String(http.StatusBadRequest, "database is empty")
	case len(cfg.Collection) == 0:
		return c.String(http.StatusBadRequest, "collection is empty")
	}
	mh.UseDbCol(cfg.Database, cfg.Collection)
	lk.Log("now: %v", cfg)
	return c.String(http.StatusOK, fmt.Sprintf("%v", cfg))
}

// @Title insert or update one entity data
// @Summary insert or update one entity data by a json file
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Param   entity  path  string true  "entity name for incoming entity data"
// @Param   data    body  string true  "entity json data for uploading" Format(binary)
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/insert/{entity} [post]
func Insert(c echo.Context) error {
	var (
		entity  = c.Param("entity")
		dataRdr = c.Request().Body
	)
	if len(entity) == 0 {
		return c.String(http.StatusBadRequest, "entity name is empty")
	}
	if dataRdr == nil {
		return c.String(http.StatusBadRequest, "payload for insert is empty")
	}

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
// @Param   data  body  string true  "json data for query" Format(binary)
// @Success 200 "OK - find successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/find [get]
func Find(c echo.Context) error {
	var (
		qryRdr = c.Request().Body
	)

	// if qryRdr == nil {
	// 	return c.String(http.StatusBadRequest, "payload for query is empty")
	// }

	results, err := mh.Find[EntityType](qryRdr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// @Router /api/entity/dump [get]
// func Dump(c echo.Context) error {
// }

// func AllEntityNames(c echo.Context) error {
// }

// func AllEntities(c echo.Context) error {
// }
