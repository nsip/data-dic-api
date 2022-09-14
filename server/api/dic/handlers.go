package dic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	"github.com/digisan/gotk/strs"
	tc "github.com/digisan/gotk/type-check"
	lk "github.com/digisan/logkit"
	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/db"
	in "github.com/nsip/data-dic-api/server/ingest"
	"github.com/tidwall/gjson"
)

// @Title insert or update one dictionary item
// @Summary insert or update one entity or collection data by json payload
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   entity body string true "entity or collection json data for uploading" Format(binary)
// @Success 200 "OK - insert or update successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/upsert [post]
func Upsert(c echo.Context) error {

	lk.Log("Enter: Post Upsert")

	var (
		dataRdr = c.Request().Body
	)

	if dataRdr != nil {
		defer dataRdr.Close()
	} else {
		return c.String(http.StatusBadRequest, "payload for insert is empty")
	}

	data, err := io.ReadAll(dataRdr)
	if err != nil {
		return c.String(http.StatusBadRequest, "entity data read error: "+err.Error())
	}

	// check payload type
	//
	inType := ""
	switch {
	case json.Unmarshal(data, &EntityType{}) != nil:
		inType = "entity"
	case json.Unmarshal(data, &CollectionType{}) != nil:
		inType = "collection"
	default:
		return c.String(http.StatusBadRequest, "invalid payload, cannot be converted to entity or collection")
	}

	js := string(data)

	// check json value type
	//
	flagHtml := false
	if strs.ContainsAny(js, "<p>", "<h1>", "<h2>", "<h3>", "<h4>", "<h5>", "<h6>") &&
		strs.ContainsAny(js, "</p>", "</h1>", "</h2>", "</h3>", "</h4>", "</h5>", "</h6>") {
		flagHtml = true
	}

	// validate payload
	//
	name := gjson.Get(js, "Entity").String()
	id := gjson.Get(js, "Metadata.Identifier").String()

	if flagHtml {
		if strings.HasPrefix(name, "<") && strings.HasSuffix(name, ">") {
			name = strs.HtmlTextContent(name)[0]
		}
	}

	switch {
	case len(name) == 0:
		return c.String(http.StatusBadRequest, "invalid payload, 'Entity' field is missing")
	case !flagHtml && !tc.IsNumeric(id):
		return c.String(http.StatusBadRequest, "invalid payload, 'Metadata.Identifier' field is invalid")
	}

	//////////////////////////////////////////////////////////

	cfg, ok := db.CfgGrp[inType]
	if !ok {
		return c.String(http.StatusBadRequest, "objType only can be [entity collection]")
	}

	mh.UseDbCol(cfg.Db, IF(flagHtml, cfg.ColHtml, cfg.ColText))

	// ingest inbound data into db(text/html), if entity already exists, replace old one
	IdOrCnt, data, err := mh.Upsert(bytes.NewReader(data), "Entity", name)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error in db upsert: "+err.Error())
	}

	// save inbound json file to local folder
	if len(data) > 0 {

		// TO inbound directory
		dir := IF(flagHtml, cfg.DirHtml, cfg.DirText)
		if err := os.WriteFile(filepath.Join(dir, name+".json"), data, os.ModePerm); err != nil {
			return c.String(http.StatusInternalServerError, "error in writing file to inbound: "+err.Error())
		}

		// TO renamed directory
		if !flagHtml {
			dir := cfg.DirExisting
			if err := os.WriteFile(filepath.Join(dir, name+".json"), data, os.ModePerm); err != nil {
				return c.String(http.StatusInternalServerError, "error in writing file to existing directory: "+err.Error())
			}
		}

		// Re ingest all, then update db(entities/collections)
		if err := in.IngestViaCmd(); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, IdOrCnt)
}

/////////////////////////////////////// FOR NEXT FRONTEND VERSION ///////////////////////////////////////

// @Title get all items
// @Summary get all entities' or collections' full content
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   itemType path  string true  "item type, only can be 'entity' or 'collection'"
// @Param   name     query string false "entity/collection 'Entity' name for query. if empty, get all"
// @Success 200 "OK - get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/items/{itemType} [get]
func Items(c echo.Context) error {

	lk.Log("Enter: Get All")

	var (
		itemType = c.Param("itemType")  // entity or collection
		name     = c.QueryParam("name") // entity or collection's "Entity" value
	)

	cfg, ok := db.CfgGrp[itemType]
	if !ok {
		return c.String(http.StatusBadRequest, "itemType can only be [entity collection]")
	}

	var (
		result any
		err    error
	)
	switch itemType {
	case "entity":
		result, err = db.Many[EntityType](cfg, name)
	case "collection":
		result, err = db.Many[CollectionType](cfg, name)
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

// @Title list all items' name
// @Summary list all entities' or collections' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   itemType path  string true  "item type, only can be 'entity' or 'collection'"
// @Param   name     query string false "entity/collection 'Entity' name for query. if empty, get all"
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/list/{itemType} [get]
func List(c echo.Context) error {

	lk.Log("Enter: Get List")

	var (
		itemType = c.Param("itemType")  // entity or collection
		name     = c.QueryParam("name") // entity or collection's "Entity" value
	)

	cfg, ok := db.CfgGrp[itemType]
	if !ok {
		return c.String(http.StatusBadRequest, "itemType can only be [entity collection]")
	}

	var (
		names []string
		err   error
	)
	switch itemType {
	case "entity":
		names, err = db.ListMany[EntityType](cfg, name)
	case "collection":
		names, err = db.ListMany[CollectionType](cfg, name)
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, names)
}

// @Title get one item
// @Summary get one entity or collection by its 'Entity' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name  query string true "Entity name"
// @Param   fuzzy query bool false "regex applies?" false
// @Success 200 "OK - got successfully"
// @Failure 400 "Fail - invalid parameters"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/one [get]
func One(c echo.Context) error {

	lk.Log("Enter: Get One")

	var (
		name     = c.QueryParam("name")
		fuzzyStr = c.QueryParam("fuzzy")
		fuzzy    = false
		err      error
	)

	if len(fuzzyStr) > 0 {
		fuzzy, err = strconv.ParseBool(fuzzyStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "'fuzzy' must be bool type, @"+err.Error())
		}
	}

	for itemType, cfg := range db.CfgGrp {
		var (
			result any
			err    error
		)
		switch itemType {
		case "entity":
			result, err = db.One[EntityType](cfg, name, fuzzy)
		case "collection":
			result, err = db.One[CollectionType](cfg, name, fuzzy)
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if result == nil {
			continue
		}
		return c.JSON(http.StatusOK, result)
	}
	return c.String(http.StatusNotFound, fmt.Sprintf(`[%v] is not existing`, name))
}

// @Title delete one item
// @Summary delete one entity or collection by its 'Entity' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name query string true "Entity name for deleting"
// @Success 200 "OK - deleted successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/one [delete]
func Delete(c echo.Context) error {

	lk.Log("Enter: Delete Delete")

	var (
		name = c.QueryParam("name")
	)
	for itemType, cfg := range db.CfgGrp {
		var (
			n   int
			err error
		)
		switch itemType {
		case "entity":
			n, err = db.Del[EntityType](cfg, name)
		case "collection":
			n, err = db.Del[CollectionType](cfg, name)
		}
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if n > 0 {
			return c.JSON(http.StatusOK, struct{ CountDeleted int }{n})
		}
	}
	return c.JSON(http.StatusOK, struct{ CountDeleted int }{0})
}

// @Title delete all items
// @Summary delete all entities or collections (dangerous)
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   itemType path string true "item type, only can be 'entity' or 'collection'"
// @Success 200 "OK - cleared successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/clear/{itemType} [delete]
func Clear(c echo.Context) error {

	lk.Log("Enter: Delete Clear")

	var (
		itemType = c.Param("itemType") // entity or collection
	)

	cfg, ok := db.CfgGrp[itemType]
	if !ok {
		return c.String(http.StatusBadRequest, "itemType can only be [entity collection]")
	}

	var (
		n   int
		err error
	)
	switch itemType {
	case "entity":
		n, err = db.Clr[EntityType](cfg)
	case "collection":
		n, err = db.Clr[CollectionType](cfg)
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, struct{ CountDeleted int }{n})
}
