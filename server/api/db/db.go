package db

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
)

func fieldStr[T any](v *T, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func One[T any](cfg Config, EntityName string, fuzzy bool) (*T, error) {
	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, EntityName)
	if fuzzy {
		sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, EntityName)
	}
	lk.Log("sFilter %v", sFilter)

	// existing db
	//
	mh.UseDbCol(cfg.Db, cfg.ColExisting)
	found, err := mh.FindOne[T](strings.NewReader(sFilter))
	if err != nil {
		return nil, err
	}
	return found, nil
}

func Many[T any](cfg Config, EntityName string) ([]*T, error) {
	var rFilter io.Reader = nil // NOT using "*strings.Reader = nil" as nil interface
	if len(EntityName) > 0 {
		sFilter := fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, EntityName)
		lk.Log("sFilter %v", sFilter)
		rFilter = strings.NewReader(sFilter)
	}

	// existing db
	//
	mh.UseDbCol(cfg.Db, cfg.ColExisting)
	found, err := mh.Find[T](rFilter)
	if err != nil {
		return nil, err
	}
	return found, nil
}

func ListMany[T any](cfg Config, EntityName string) ([]string, error) {
	found, err := Many[T](cfg, EntityName)
	if err != nil {
		return nil, err
	}
	return FilterMap(found, nil, func(i int, e *T) string { return fieldStr(e, "Entity") }), nil
}

func Del[T any](cfg Config, EntityName string) (int, error) {

	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, EntityName)
	lk.Log("sFilter %v", sFilter)

	// existing db
	//
	mh.UseDbCol(cfg.Db, cfg.ColExisting)
	nExisting, _, err := mh.DeleteOne[T](strings.NewReader(sFilter)) // MUST re-create reader
	if err != nil {
		return 0, err
	}

	// inbound db, text
	//
	mh.UseDbCol(cfg.Db, cfg.ColText)
	nText, _, err := mh.DeleteOne[T](strings.NewReader(sFilter)) // MUST re-create reader
	if err != nil {
		return 0, err
	}

	// inbound db, html (may change in future)
	//
	mh.UseDbCol(cfg.Db, cfg.ColHtml)
	sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)>?%v<?"}}`, EntityName)
	nHtml, _, err := mh.DeleteOne[T](strings.NewReader(sFilter)) // MUST re-create reader
	if err != nil {
		return 0, err
	}

	return Max(nExisting, nText, nHtml), nil
}

func Clr[T any](cfg Config) (int, error) {
	names, err := ListMany[T](cfg, "")
	if err != nil {
		return 0, err
	}
	cnt := 0
	for _, name := range names {
		n, err := Del[T](cfg, name)
		if err != nil {
			return 0, err
		}
		cnt += n
	}
	return cnt, nil
}
