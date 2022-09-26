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

var (
	mListCache = make(map[string][]string)
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
// @Router /api/dictionary/auth/upsert [post]
// @Security ApiKeyAuth
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
	case json.Unmarshal(data, &EntityType{}) == nil:
		inType = "entity"
	case json.Unmarshal(data, &CollectionType{}) == nil:
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
		if err := in.IngestViaCmd(false); err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return c.JSON(http.StatusOK, IdOrCnt)
}

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
// @Router /api/dictionary/pub/items/{itemType} [get]
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
// @Router /api/dictionary/pub/list/{itemType} [get]
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

	// cache list
	mListCache[itemType] = names

	return c.JSON(http.StatusOK, names)
}

// @Title get one item
// @Summary get one entity or collection by its 'Entity' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name  query string  true "Entity name"
// @Param   fuzzy query boolean false "regex applies?" false
// @Success 200 "OK - got successfully"
// @Failure 400 "Fail - invalid parameters"
// @Failure 404 "Fail - not found"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/one [get]
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
			result any = nil
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
		if !tc.IsNil(result) {
			return c.JSON(http.StatusOK, result)
		}
	}

	lk.Warn("Not Found: [%v]", name)

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
// @Router /api/dictionary/auth/one [delete]
// @Security ApiKeyAuth
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
// @Router /api/dictionary/auth/clear/{itemType} [delete]
// @Security ApiKeyAuth
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

// @Title check item's kind by its name
// @Summary check item's kind ('entity' or 'collection') by its name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name query string true "Entity name for checking kind"
// @Success 200 "OK - got kind ('entity' or 'collection') successfully"
// @Failure 404 "Fail - neither 'entity' nor 'collection'"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/kind [get]
func CheckItemKind(c echo.Context) error {

	lk.Log("Enter: CheckItemKind")

	var (
		name = c.QueryParam("name")
	)

	if len(mListCache) == 0 {
		return c.String(http.StatusInternalServerError, "list cache hasn't been loaded")
	}
	for kind, list := range mListCache {
		if len(list) == 0 {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("[%s] list hasn't been loaded", kind))
		}
	}

	switch {
	case strs.IsIn(true, true, name, mListCache["entity"]...):
		return c.JSON(http.StatusOK, "entity")

	case strs.IsIn(true, true, name, mListCache["collection"]...):
		return c.JSON(http.StatusOK, "collection")

	default:
		return c.String(http.StatusNotFound, "unknown")
	}
}

//////////////////////////////////////////////////////////////////

// @Title get related entities of a collection
// @Summary get related entities' name of a collection
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   colname query string true "collection name"
// @Success 200 "OK - got collection content successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/colentities [get]
func ColEntities(c echo.Context) error {

	lk.Log("Enter: ColEntities")

	var (
		name = c.QueryParam("colname")
	)
	rt, err := db.ColEntities(db.CfgGrp["collection"], name) // only use cfg's dbName, colName is fixed.
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rt)
}

// @Title get class info of an entity
// @Summary get class info (derived path & children) of an entity
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   entname query string true "entity name"
// @Success 200 "OK - got entity class info successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/entclasses [get]
func EntClasses(c echo.Context) error {

	lk.Log("Enter: EntClasses")

	var (
		name = c.QueryParam("entname")
	)
	derived, children, err := db.EntClasses(db.CfgGrp["entity"], name) // only use cfg's dbName, colName is fixed.
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		DerivedPath []string
		Children    []string
	}{derived, children})
}

// @Title get list of item's name by searching
// @Summary get list of entity's & collection's name by searching. If not given, return all
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   aim        query string  false "search content from whole dictionary"
// @Param   ignorecase query boolean false "case insensitive ?"
// @Success 200 "OK - got list of found item's name successfully"
// @Failure 400 "Fail - invalid parameters"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/search [get]
func FullTextSearch(c echo.Context) error {

	lk.Log("Enter: FullTextSearch")

	var (
		aim        = c.QueryParam("aim")
		ignorecase = c.QueryParam("ignorecase")
	)

	// if aim is empty, return list of all items
	if len(aim) == 0 {
		return c.JSON(http.StatusOK, struct {
			Entities    []string
			Collections []string
		}{mListCache["entity"], mListCache["collection"]})
	}

	// aim is not empty, do real full text search
	flagIgnCase, err := strconv.ParseBool(ignorecase)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	entities, collections, err := db.FullTextSearch(db.CfgGrp["entity"], aim, flagIgnCase) // only use cfg's dbName, colName is fixed.
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		Entities    []string
		Collections []string
	}{entities, collections})
}
