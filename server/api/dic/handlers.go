package dic

import (
	"bytes"
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
	u "github.com/digisan/user-mgr/user"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/nsip/data-dic-api/server/api/db"
	"github.com/tidwall/gjson"
)

// @Title   submit one dictionary item
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
		userTkn = c.Get("user").(*jwt.Token)     //
		claims  = userTkn.Claims.(*u.UserClaims) //
		author  = claims.UName                   // author
		dataRdr = c.Request().Body               // data
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

	//
	// check payload kind
	//
	inKind := db.ItemKind(data)
	if len(inKind) == 0 {
		return c.String(http.StatusBadRequest, "invalid payload, cannot be converted to entity or collection")
	}

	//
	// payload json string
	//
	js := string(data)

	//
	// check json value type
	//
	flagHtml := false
	if strs.ContainsAny(js, "<p>", "<h1>", "<h2>", "<h3>", "<h4>", "<h5>", "<h6>") &&
		strs.ContainsAny(js, "</p>", "</h1>", "</h2>", "</h3>", "</h4>", "</h5>", "</h6>") {
		flagHtml = true
	}

	//
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

	cfg, ok := db.CfgGrp[inKind]
	if !ok {
		return c.String(http.StatusBadRequest, "item kind can only be [entity collection]")
	}

	mh.UseDbCol(db.DATABASE, string(IF(flagHtml, cfg.DbColHtml, cfg.DbColText)))

	//
	// ingest inbound data into db(text/html), if entity already exists, replace old one
	//
	IdOrCnt, data, err := mh.Upsert(bytes.NewReader(data), "Entity", name)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, "error in db upsert: "+err.Error())
	}

	//
	// save inbound json file to local folder
	//
	if len(data) > 0 {
		// to inbound directory
		dir := IF(flagHtml, cfg.DirHtml, cfg.DirText)
		if err := os.WriteFile(filepath.Join(dir, name+".json"), data, os.ModePerm); err != nil {
			lk.WarnOnErr("%v", err)
			return c.String(http.StatusInternalServerError, "error in writing file to inbound: "+err.Error())
		}
	}

	//////////////////////////////////////////////////////////

	//
	// record item(text) into author db collection
	//
	if !flagHtml {
		if _, err := db.RecordAction(author, db.Submit, name, inKind); err != nil {
			lk.WarnOnErr("%v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, IdOrCnt)
}

// @Title   get all items
// @Summary get all entities' or collections' full content
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   kind  path  string true  "item type, can only be [entity collection]"
// @Param   name  query string false "entity/collection 'Entity' name for query. if empty, get all"
// @Param   dbcol query string true  "from which db collection? [existing, text, html]"
// @Success 200 "OK - get successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/items/{kind} [get]
func Items(c echo.Context) error {

	lk.Log("Enter: Get All")

	var (
		kind  = c.Param("kind")       // entity or collection
		name  = c.QueryParam("name")  // entity or collection's "Entity" value
		dbcol = c.QueryParam("dbcol") // existing, text, html
	)

	cfg, ok := db.CfgGrp[kind]
	if !ok {
		return c.String(http.StatusBadRequest, "'kind' can only be [entity collection]")
	}

	var (
		result any
		err    error
	)
	switch kind {
	case "entity":
		result, err = db.Many[db.EntityType](cfg, db.DbColType(dbcol), name)
	case "collection":
		result, err = db.Many[db.CollectionType](cfg, db.DbColType(dbcol), name)
	}
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

// @Title   list all items' name
// @Summary list all entities' or collections' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   kind  path  string true  "item type, can only be [entity collection]"
// @Param   name  query string false "entity/collection 'Entity' name for query. if empty, get all"
// @Param   dbcol query string true  "from which db collection? [existing, text, html]"
// @Success 200 "OK - list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/list/{kind} [get]
func List(c echo.Context) error {
	mtx.Lock()
	defer mtx.Unlock()

	lk.Log("Enter: Get List")

	var (
		kind  = c.Param("kind")       // entity or collection
		name  = c.QueryParam("name")  // entity or collection's "Entity" value
		dbcol = c.QueryParam("dbcol") // existing, text, html
	)

	cfg, ok := db.CfgGrp[kind]
	if !ok {
		return c.String(http.StatusBadRequest, "'kind' can only be [entity collection]")
	}

	var (
		names []string
		err   error
	)
	switch kind {
	case "entity":
		names, err = db.ListMany[db.EntityType](cfg, db.DbColType(dbcol), name)
	case "collection":
		names, err = db.ListMany[db.CollectionType](cfg, db.DbColType(dbcol), name)
	}
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// cache list
	mListCache[dbcol][kind] = names

	return c.JSON(http.StatusOK, names)
}

// @Title   get one item
// @Summary get one entity or collection by its 'Entity' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name  query string  true "Entity name"
// @Param   fuzzy query boolean false "regex applies?" false
// @Param   dbcol query string true  "from which db collection? [existing, text, html]"
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
		dbcol    = c.QueryParam("dbcol") // existing, text, html
		fuzzy    = false
		err      error
	)

	if len(fuzzyStr) > 0 {
		fuzzy, err = strconv.ParseBool(fuzzyStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "'fuzzy' must be bool type, @"+err.Error())
		}
	}

	for kind, cfg := range db.CfgGrp {
		var (
			result any = nil
			err    error
		)
		switch kind {
		case "entity":
			result, err = db.One[db.EntityType](cfg, db.DbColType(dbcol), name, fuzzy)
		case "collection":
			result, err = db.One[db.CollectionType](cfg, db.DbColType(dbcol), name, fuzzy)
		}
		if err != nil {
			lk.WarnOnErr("%v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if !tc.IsNil(result) {
			return c.JSON(http.StatusOK, result)
		}
	}

	lk.Warn("Not Found: [%v]", name)

	return c.String(http.StatusNotFound, fmt.Sprintf(`[%v] is not existing`, name))
}

// @Title   check item existing status
// @Summary check whether one item exists according to its 'Entity' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name query string true "Entity name"
// @Param   dbcol query string true  "from which db collection? [existing, text, html]"
// @Success 200 "OK - got successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/pub/exists [get]
func Exists(c echo.Context) error {

	lk.Log("Enter: Get Exists")

	var (
		name  = c.QueryParam("name")
		dbcol = c.QueryParam("dbcol") // existing, text, html
	)
	ok, err := db.Exists(db.DbColType(dbcol), name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ok)
}

// @Title   delete one item
// @Summary delete one entity or collection by its 'Entity' name
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name  query string true "Entity name for deleting"
// @Param   dbcol query string true "from which db collection? [existing, text, html]"
// @Success 200 "OK - deleted successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/auth/one [delete]
// @Security ApiKeyAuth
func Delete(c echo.Context) error {

	lk.Log("Enter: Delete Delete")

	var (
		name  = c.QueryParam("name")
		dbcol = c.QueryParam("dbcol") // existing, text, html
	)
	for kind, cfg := range db.CfgGrp {
		var (
			n   int
			err error
		)
		switch kind {
		case "entity":
			n, err = db.Del[db.EntityType](cfg, db.DbColType(dbcol), name)
		case "collection":
			n, err = db.Del[db.CollectionType](cfg, db.DbColType(dbcol), name)
		}
		if err != nil {
			lk.WarnOnErr("%v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		if n > 0 {
			return c.JSON(http.StatusOK, struct{ CountDeleted int }{n})
		}
	}
	return c.JSON(http.StatusOK, struct{ CountDeleted int }{0})
}

// @Title   delete all items
// @Summary delete all entities or collections (dangerous)
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   kind  path  string true "item type, can only be [entity collection]"
// @Param   dbcol query string true "which db collection? [existing, text, html]"
// @Success 200 "OK - cleared successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/auth/clear/{kind} [delete]
// @Security ApiKeyAuth
func Clear(c echo.Context) error {

	lk.Log("Enter: Delete Clear")

	var (
		kind  = c.Param("kind")       // entity or collection
		dbcol = c.QueryParam("dbcol") // existing, text, html
	)

	cfg, ok := db.CfgGrp[kind]
	if !ok {
		return c.String(http.StatusBadRequest, "'kind' can only be [entity collection]")
	}

	var (
		n   int
		err error
	)
	switch kind {
	case "entity":
		n, err = db.Clr[db.EntityType](cfg, db.DbColType(dbcol))
	case "collection":
		n, err = db.Clr[db.CollectionType](cfg, db.DbColType(dbcol))
	}
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, struct{ CountDeleted int }{n})
}

// @Title   check item's kind by its name
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
func ItemKind(c echo.Context) error {
	mtx.Lock()
	defer mtx.Unlock()

	lk.Log("Enter: CheckItemKind")

	var (
		name = c.QueryParam("name")
	)

	if len(mListCache["existing"]) == 0 {
		return c.String(http.StatusInternalServerError, "list cache hasn't been loaded")
	}

	lsEntity, okEntity := mListCache["existing"]["entity"]
	lsCollection, okCollection := mListCache["existing"]["collection"]

	switch {
	case okEntity && strs.IsIn(true, true, name, lsEntity...):
		return c.JSON(http.StatusOK, "entity")

	case okCollection && strs.IsIn(true, true, name, lsCollection...):
		return c.JSON(http.StatusOK, "collection")

	default:
		return c.String(http.StatusNotFound, "unknown")
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////

// @Title   get related entities of a collection
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
	rt, err := db.GetColEntities(name) // only use cfg's dbName, colName is fixed.
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, rt)
}

// @Title   get class info of an entity
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
	derived, children, err := db.GetEntClasses(name) // only use cfg's dbName, colName is fixed.
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		DerivedPath []string
		Children    []string
	}{derived, children})
}

// @Title   get list of item's name by searching
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
	mtx.Lock()
	defer mtx.Unlock()

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
		}{mListCache["existing"]["entity"], mListCache["existing"]["collection"]})
	}

	// aim is not empty, do real full text search
	flagIgnCase, err := strconv.ParseBool(ignorecase)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	entities, collections, err := db.FullTextSearch(aim, flagIgnCase)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		Entities    []string
		Collections []string
	}{entities, collections})
}

// @Title   approve one dictionary item
// @Summary approve one dictionary candidate item
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name query string true "entity/collection 'Entity' name for approval"
// @Param   kind query string true "item type, can only be [entity collection]"
// @Success 200 "OK - approve successfully"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 404 "Fail - couldn't find item to approve"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/auth/approve [put]
// @Security ApiKeyAuth
func Approve(c echo.Context) error {

	lk.Log("Enter: Put Approve")

	var (
		userTkn  = c.Get("user").(*jwt.Token)     //
		claims   = userTkn.Claims.(*u.UserClaims) //
		approver = claims.UName                   // approver
		name     = c.QueryParam("name")           // item name
		kind     = c.QueryParam("kind")           // item kind
	)

	cfg, ok := db.CfgGrp[kind]
	if !ok {
		return c.String(http.StatusBadRequest, "kind can only be [entity collection]")
	}

	//
	// 1. delete from inbound text db collection. KEEP html for future format query
	//
	mh.UseDbCol(db.DATABASE, string(cfg.DbColText))

	n, _, err := mh.DeleteOneAt[any]("Entity", name)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if n != 1 {
		lk.Warn("deleted N: %v", n)
		return c.String(http.StatusNotFound, fmt.Sprintf("couldn't find [%v], approved nothing", name))
	}

	//
	// 2. record into existing db collection
	//
	mh.UseDbCol(db.DATABASE, string(cfg.DbColExisting))

	var (
		from = filepath.Join(cfg.DirText, name+".json")
		to   = filepath.Join(cfg.DirExisting, name+".json")
	)

	file, err := os.Open(from)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	_, _, err = mh.Upsert(file, "Entity", name)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	//
	// 3. move file from text to existing
	//
	if err := os.Rename(from, to); err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	//
	// 4. record into approver db collection
	//
	if _, err := db.RecordAction(approver, db.Approve, name, kind); err != nil {
		lk.WarnOnErr("%v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("[%v] is approved by [%v]", name, approver))
}

// @Title   toggle subscribe one dictionary item
// @Summary toggle subscribe one dictionary item
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   name query string true "entity/collection 'Entity' name for approval"
// @Param   kind query string true "item type, can only be [entity collection]"
// @Success 200 "OK - subscribe/unsubscribe successfully. true: subscribed now; false: unsubscribed now"
// @Failure 400 "Fail - invalid parameters or request body"
// @Failure 404 "Fail - couldn't find item to subscribe"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/auth/subscribe [put]
// @Security ApiKeyAuth
func ToggleSubscribe(c echo.Context) error {

	lk.Log("Enter: Put Subscribe")

	var (
		userTkn = c.Get("user").(*jwt.Token)     //
		claims  = userTkn.Claims.(*u.UserClaims) //
		user    = claims.UName                   // user
		name    = c.QueryParam("name")           // item name
		kind    = c.QueryParam("kind")           // item kind
	)

	// validate 'kind'
	if _, ok := db.CfgGrp[kind]; !ok {
		return c.String(http.StatusBadRequest, "[kind] can only be [entity collection]")
	}

	// validate 'name'
	ok, err := db.Exists("existing", name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if !ok {
		return c.String(http.StatusNotFound, fmt.Sprintf("[%v] is not existing, cannot subscribe", name))
	}

	// check subscription status
	ls, err := db.ListActionRecord(user, db.DbColType("subscribe"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if NotIn(name, ls...) {

		// subscribe action
		_, err := db.RecordAction(user, db.Subscribe, name, kind)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, true)

	} else {

		// unsubscribe action
		_, err := db.RemoveAction(user, db.Subscribe, name)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, false)
	}
}

// @Title   list action items
// @Summary list recorded action items from [submit, approve, subscribe]
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   action path string true "which action what to list its item record"
// @Success 200 "OK - get list successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/auth/list/{action} [get]
// @Security ApiKeyAuth
func ListAction(c echo.Context) error {

	lk.Log("Enter: Get ListAction")

	var (
		userTkn = c.Get("user").(*jwt.Token)     //
		claims  = userTkn.Claims.(*u.UserClaims) //
		user    = claims.UName                   // user
		action  = c.Param("action")              // action type: submit, approve, subscribe
	)

	ls, err := db.ListActionRecord(user, db.DbColType(action))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ls)
}

// @Title   check action item is existing
// @Summary check recorded action item existing status from [submit, approve, subscribe]
// @Description
// @Tags    Dictionary
// @Accept  json
// @Produce json
// @Param   action path  string true "which action what to list its item record"
// @Param   name   query string true "entity/collection 'Entity' name for checking existing status"
// @Success 200 "OK - get existing status successfully"
// @Failure 500 "Fail - internal error"
// @Router /api/dictionary/auth/check/{action} [get]
// @Security ApiKeyAuth
func ActionRecordExists(c echo.Context) error {

	lk.Log("Enter: Get ActionRecordExists")

	var (
		userTkn = c.Get("user").(*jwt.Token)     //
		claims  = userTkn.Claims.(*u.UserClaims) //
		user    = claims.UName                   // user
		action  = c.Param("action")              // action: submit, approve, subscribe
		name    = c.QueryParam("name")           // item name
	)

	ls, err := db.ListActionRecord(user, db.DbColType(action))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, In(name, ls...))
}
