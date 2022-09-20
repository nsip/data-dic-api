package db

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
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

func One[T any](cfg Config, EntName string, fuzzy bool) (*T, error) {
	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, EntName)
	if fuzzy {
		sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, EntName)
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

func Many[T any](cfg Config, EntName string) ([]*T, error) {
	var rFilter io.Reader = nil // NOT using "*strings.Reader = nil" as nil interface
	if len(EntName) > 0 {
		sFilter := fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, EntName)
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

func ListMany[T any](cfg Config, EntName string) ([]string, error) {
	found, err := Many[T](cfg, EntName)
	if err != nil {
		return nil, err
	}
	return FilterMap(found, nil, func(i int, e *T) string { return fieldStr(e, "Entity") }), nil
}

func Del[T any](cfg Config, EntName string) (int, error) {

	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, EntName)
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
	sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)>?%v<?"}}`, EntName)
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

/////////////////////////////////////////////////////////////

func ColEntities(cfg Config, ColName string) ([]string, error) {

	mh.UseDbCol(cfg.Db, "colentities") // fixed collection name

	whole, err := mh.FindOne[map[string]any](nil) // only one
	if err != nil {
		return nil, err
	}
	rtAny, ok := (*whole)[ColName] // primitive.A
	if !ok {
		return []string{}, nil
	}

	rtStr, err := mh.CvtA[string](rtAny)
	if err != nil {
		return nil, err
	}
	return rtStr, nil
}

func EntClasses(cfg Config, EntName string) ([]string, []string, error) {

	mh.UseDbCol(cfg.Db, "class") // fixed collection name

	whole, err := mh.FindOne[map[string]any](nil) // only one
	if err != nil {
		return nil, nil, err
	}
	rt, ok := (*whole)[EntName] // primitive.M
	if !ok {
		return []string{}, []string{}, nil
	}

	c, err := mh.CvtM[struct {
		Branch   string
		Children []string
	}](rt)
	if err != nil {
		return nil, nil, err
	}
	return strings.Split(c.Branch, "--"), c.Children, nil
}

func FullTextSearch(cfg Config, aim string, insensitive bool) ([]string, []string, error) {

	mh.UseDbCol(cfg.Db, "pathval") // fixed collection name

	var (
		entities    = []string{}
		collections = []string{}
		err         error
	)

	all, err := mh.Find[map[string]any](nil) // need to iterate all
	if err != nil {
		return nil, nil, err
	}

	prefix := IF(insensitive, "(?i)", "")
	sep := `[\s\\/-:]+`
	r0 := regexp.MustCompile(fmt.Sprintf(`%s^%s$`, prefix, aim))             // whole
	r1 := regexp.MustCompile(fmt.Sprintf(`%s^%s%s`, prefix, aim, sep))       // start
	r2 := regexp.MustCompile(fmt.Sprintf(`%s%s%s%s`, prefix, sep, aim, sep)) // middle
	r3 := regexp.MustCompile(fmt.Sprintf(`%s%s%s$`, prefix, sep, aim))       // end

	for _, one := range all {
		itemName := (*one)["Entity"].(string)

		flagCol := false
		switch (*one)["Metadata[dot]Type"].(string) {
		case "Collection", "collection":
			flagCol = true
		}

		for _, val := range *one {
			valstr := strings.TrimSpace(fmt.Sprint(val))
			if r0.MatchString(valstr) || r1.MatchString(valstr) || r2.MatchString(valstr) || r3.MatchString(valstr) {
				if flagCol {
					collections = append(collections, itemName)
				} else {
					entities = append(entities, itemName)
				}
				break
			}
		}
	}
	return entities, collections, err
}
