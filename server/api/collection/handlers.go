package collection

import (
	"bytes"
	"fmt"
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

// @Title insert or update one collection data
// @Summary insert or update one collection data by a json file
// @Description
// @Tags    Collection
// @Accept  json
// @Produce json
// @Param   valType     path string true "[text] value or [html] value in json payload"
// @Param   collection  body string true "collection json data for uploading" Format(binary)
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/collection/upsert/{valType} [post]
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
		return c.String(http.StatusBadRequest, "collection data read error: "+err.Error())
	}

	js := string(data)

	// validate payload
	collectionName := gjson.Get(js, "Entity").String()
	collectionId := gjson.Get(js, "Metadata.Identifier").String()

	if flagHtml {
		if strings.HasPrefix(collectionName, "<") && strings.HasSuffix(collectionName, ">") {
			collectionName = strs.HtmlTextContent(collectionName)[0]
		}
	}

	switch {
	case len(collectionName) == 0:
		return c.String(http.StatusBadRequest, "invalid collection json, 'Entity' field is missing")
	case !flagHtml && !tc.IsNumeric(collectionId):
		return c.String(http.StatusBadRequest, "invalid collection json, 'Metadata.Identifier' field is invalid")
	}

	// ingest inbound data into db, if collection already exists, replace old one
	IdOrCnt, data, err := mh.Upsert(bytes.NewReader(data), "Entity", collectionName)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in db upsert: "+err.Error())
	}

	// save inbound json file to local folder
	if len(data) > 0 {

		// TO inbound
		dir := IF(flagHtml, cfg.dirHtml, cfg.dirText)
		if err := os.WriteFile(filepath.Join(dir, collectionName+".json"), data, os.ModePerm); err != nil {
			return c.String(http.StatusInternalServerError, "error in writing file to inbound: "+err.Error())
		}

		// TO renamed
		if !flagHtml {
			dir := cfg.dirExisting
			if err := os.WriteFile(filepath.Join(dir, collectionName+".json"), data, os.ModePerm); err != nil {
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

// @Title get all collections
// @Summary get all collections' full content
// @Description
// @Tags    Collection
// @Accept  json
// @Produce json
// @Success 200 "OK - get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/collection/collections [get]
func AllCollections(c echo.Context) error {
	collections, err := allCollections()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, collections)
}

// @Title list all collection names
// @Summary list all collection names
// @Description
// @Tags    Collection
// @Accept  json
// @Produce json
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/collection/names [get]
func AllCollectionNames(c echo.Context) error {
	names, err := allCollectionNames()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, names)
}

// @Title delete one collection
// @Summary delete one collection by its name
// @Description
// @Tags    Collection
// @Accept  json
// @Produce json
// @Param   collection query string true "collection's entity name for collection deleting"
// @Success 200 "OK - deleted successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/collection/name [delete]
func Delete(c echo.Context) error {
	var (
		collection = c.QueryParam("collection")
	)
	n, _, err := mh.DeleteOne[CollectionType](strings.NewReader(fmt.Sprintf(`"Entity": "%v"`, collection)))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct{ CountDeleted int }{n})
}

// @Title delete all collections
// @Summary delete all collections (dangerous)
// @Description
// @Tags    Collection
// @Accept  json
// @Produce json
// @Success 200 "OK - deleted successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/collection/clear_all [delete]
func ClearAll(c echo.Context) error {
	names, err := allCollectionNames()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	total := 0
	for _, name := range names {
		n, _, err := mh.Delete[CollectionType](strings.NewReader(fmt.Sprintf(`"Entity": "%v"`, name)))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		total += n
	}
	return c.JSON(http.StatusOK, struct{ CountDeleted int }{total})
}
