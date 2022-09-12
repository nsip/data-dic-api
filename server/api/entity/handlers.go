package entity

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	tc "github.com/digisan/gotk/type-check"
	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
	in "github.com/nsip/data-dic-api/server/ingest"
	"github.com/tidwall/gjson"
)

// @Title find entity json
// @Summary find entities json content by pass a json query string via payload
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Param   data body string true "json data for query" Format(binary)
// @Success 200 "OK - find successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/find [get]
// func Find(c echo.Context) error {

// 	lk.Log("Enter: Find")

// 	var (
// 		qryRdr = c.Request().Body
// 	)

// 	if qryRdr != nil {
// 		defer qryRdr.Close()
// 	} else {
// 		return c.String(http.StatusBadRequest, "payload for query is empty")
// 	}

// 	results, err := mh.Find[EntityType](qryRdr)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, results)
// }

// @Title insert or update one entity data
// @Summary insert or update one entity data by a json file
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Param   valType path string true "[text] value or [html] value in json payload"
// @Param   entity  body string true "entity json data for uploading" Format(binary)
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/upsert/{valType} [post]
func Upsert(c echo.Context) error {

	lk.Log("Enter: Upsert")

	var (
		plValType = c.Param("valType") // text or html
		dataRdr   = c.Request().Body
		flagHtml  = false
	)

	switch plValType {
	case "text", "Text", "TEXT":
		flagHtml = false
	case "html", "Html", "HTML":
		flagHtml = true
	default:
		return c.String(http.StatusBadRequest, "valType only can be [text html]")
	}

	mh.UseDbCol(cfg.db, IF(flagHtml, cfg.colHtml, cfg.colText))

	if dataRdr != nil {
		defer dataRdr.Close()
	} else {
		return c.String(http.StatusBadRequest, "payload for insert is empty")
	}

	data, err := io.ReadAll(dataRdr)
	if err != nil {
		return c.String(http.StatusBadRequest, "entity data read error: "+err.Error())
	}

	js := string(data)

	// validate payload
	entityName := gjson.Get(js, "Entity").String()
	entityId := gjson.Get(js, "Metadata.Identifier").String()

	if flagHtml {
		if strings.HasPrefix(entityName, "<") && strings.HasSuffix(entityName, ">") {
			entityName = strs.HtmlTextContent(entityName)[0]
		}
	}

	switch {
	case len(entityName) == 0:
		return c.String(http.StatusBadRequest, "invalid entity json, 'Entity' field is missing")
	case !flagHtml && !tc.IsNumeric(entityId):
		return c.String(http.StatusBadRequest, "invalid entity json, 'Metadata.Identifier' field is invalid")
	}

	// ingest inbound data into db, if entity already exists, replace old one
	IdOrCnt, data, err := mh.Upsert(bytes.NewReader(data), "Entity", entityName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in db upsert: "+err.Error())
	}

	// save inbound json file to local folder
	if len(data) > 0 {

		// TO inbound
		dir := IF(flagHtml, cfg.dirHtml, cfg.dirText)
		if err := os.WriteFile(filepath.Join(dir, entityName+".json"), data, os.ModePerm); err != nil {
			return c.String(http.StatusInternalServerError, "error in writing file to inbound: "+err.Error())
		}

		// TO renamed
		if !flagHtml {
			dir := cfg.dirExisting
			if err := os.WriteFile(filepath.Join(dir, entityName+".json"), data, os.ModePerm); err != nil {
				return c.String(http.StatusInternalServerError, "error in writing file to existing directory: "+err.Error())
			}
		}

		// Re ingest all
		if err := in.IngestViaCmd(); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, IdOrCnt)
}

/////////////////////////////////////// FOR NEXT FRONTEND VERSION ///////////////////////////////////////

// @Title get all entities
// @Summary get all entities's full content
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Success 200 "OK - get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/entities [get]
func AllEntities(c echo.Context) error {
	entities, err := allEntities()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, entities)
}

// @Title list all entity names
// @Summary list all entity names
// @Description
// @Tags    Entity
// @Accept  json
// @Produce json
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/list_names [get]
func AllEntityNames(c echo.Context) error {
	entities, err := allEntities()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	names := FilterMap(entities, nil, func(i int, e *EntityType) string { return e.Entity })
	return c.JSON(http.StatusOK, names)
}
