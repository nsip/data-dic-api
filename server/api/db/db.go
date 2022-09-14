package db

import (
	"fmt"
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

func All[T any](cfg Config) ([]*T, error) {

	// existing db
	//
	mh.UseDbCol(cfg.Db, cfg.ColExisting)
	existing, err := mh.Find[T](nil)
	if err != nil {
		return nil, err
	}
	return existing, nil
}

func List[T any](cfg Config) ([]string, error) {
	all, err := All[T](cfg)
	if err != nil {
		return nil, err
	}
	return FilterMap(all, nil, func(i int, e *T) string { return fieldStr(e, "Entity") }), nil
}

func One[T any](cfg Config, entityName string) (*T, error) {

	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, entityName)
	lk.Log("sFilter %v", sFilter)

	// existing db
	//
	mh.UseDbCol(cfg.Db, cfg.ColExisting)
	existing, err := mh.FindOne[T](strings.NewReader(sFilter))
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return nil, nil
}

func Del[T any](cfg Config, entityName string) (int, error) {

	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, entityName)
	lk.Log("sFilter %v", sFilter)

	// existing db
	//
	mh.UseDbCol(cfg.Db, cfg.ColExisting)
	nExisting, _, err := mh.DeleteOne[T](strings.NewReader(sFilter))
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

	// inbound db, html ( need to implement a new rFilter )
	//
	// mh.UseDbCol(cfg.Db, cfg.ColHtml)

	///////////////////////////////////////

	lk.WarnOnErrWhen(nText > nExisting, "%v", fmt.Errorf("nExisting must NOT less than nText"))

	return nText, nil
}

func Clr[T any](cfg Config) (int, error) {
	names, err := List[T](cfg)
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
