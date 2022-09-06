package entity

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

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

	lk.Log("Enter: UseDbCol")

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
// @Param   entityName path string true "entity name for incoming entity data"
// @Param   entityData body string true "entity json data for uploading" Format(binary)
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/upsert/{entityName} [post]
func Upsert(c echo.Context) error {

	lk.Log("Enter: Upsert")

	var (
		entityName = c.Param("entityName")
		dataRdr    = c.Request().Body
	)
	if len(entityName) == 0 {
		return c.String(http.StatusBadRequest, "entity name is empty")
	}
	if dataRdr != nil {
		defer dataRdr.Close()
	} else {
		return c.String(http.StatusBadRequest, "payload for insert is empty")
	}

	// ingest inbound data into db, if entity already exists, replace old one
	IdOrCnt, data, err := mh.Upsert(dataRdr, "Entity", entityName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in db insert"+err.Error())
	}

	// ingest inbound data into db
	// id, data, err := mh.Insert(dataRdr)
	// if err != nil {
	// 	return c.String(http.StatusInternalServerError, "error in db insert"+err.Error())
	// }

	if len(data) > 0 {
		// save inbound json file to local folder
		if err := os.WriteFile(filepath.Join(dataFolder, entityName), data, os.ModePerm); err != nil {
			return c.String(http.StatusInternalServerError, "error in writing file"+err.Error())
		}
	}

	return c.JSON(http.StatusOK, IdOrCnt)
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

	if qryRdr != nil {
		defer qryRdr.Close()
	} else {
		return c.String(http.StatusBadRequest, "payload for query is empty")
	}

	results, err := mh.Find[EntityType](qryRdr)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, results)
}

// func AllEntityNames(c echo.Context) error {
// }

// func AllEntities(c echo.Context) error {
// }
