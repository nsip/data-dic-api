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
	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

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

	col := IF(flagHtml, cfg.colHtml, cfg.colText)
	mh.UseDbCol(cfg.db, col)

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
	entityName := gjson.Get(js, "Entity").String()
	if len(entityName) == 0 {
		return c.String(http.StatusBadRequest, "invalid entity json file, 'Entity' field is missing")
	}

	if flagHtml {
		if strings.HasPrefix(entityName, "<") && strings.HasSuffix(entityName, ">") {
			entityName = strs.HtmlTextContent(entityName)[0]
		}
	}

	// ingest inbound data into db, if entity already exists, replace old one
	IdOrCnt, data, err := mh.Upsert(bytes.NewReader(data), "Entity", entityName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in db upsert: "+err.Error())
	}

	// save inbound json file to local folder
	if len(data) > 0 {
		dir := IF(flagHtml, cfg.dirHtml, cfg.dirText)
		if err := os.WriteFile(filepath.Join(dir, entityName+".json"), data, os.ModePerm); err != nil {
			return c.String(http.StatusInternalServerError, "error in writing file: "+err.Error())
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
// @Param   data body string true "json data for query" Format(binary)
// @Success 200 "OK - find successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/entity/find [get]
func Find(c echo.Context) error {

	lk.Log("Enter: Find")

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
