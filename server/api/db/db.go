package db

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
)

func fieldStr[T any](v *T, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

// from: existing, text, html
func One[T any](cfg ItemConfig, from DbColType, itemName string, fuzzy bool) (*T, error) {
	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, itemName)
	if fuzzy {
		sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, itemName)
	}
	lk.Log("sFilter %v", sFilter)

	COL := cfg.DbColVal(from)
	if len(COL) == 0 {
		return nil, fmt.Errorf("from DbCol can only be [existing, text, html]")
	}

	mh.UseDbCol(DATABASE, COL)
	found, err := mh.FindOne[T](strings.NewReader(sFilter))
	if err != nil {
		return nil, err
	}
	return found, nil
}

// from: existing, text, html
func Many[T any](cfg ItemConfig, from DbColType, itemName string) ([]*T, error) {
	var rFilter io.Reader = nil // NOT using "*strings.Reader = nil" as nil interface
	if len(itemName) > 0 {
		sFilter := fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, itemName)
		lk.Log("sFilter %v", sFilter)
		rFilter = strings.NewReader(sFilter)
	}

	COL := cfg.DbColVal(from)
	if len(COL) == 0 {
		return nil, fmt.Errorf("from ItemDbCol can only be [existing, text, html]")
	}

	mh.UseDbCol(DATABASE, COL)
	found, err := mh.Find[T](rFilter)
	if err != nil {
		return nil, err
	}
	return found, nil
}

// from: existing, text, html
func ListMany[T any](cfg ItemConfig, from DbColType, itemName string) ([]string, error) {
	found, err := Many[T](cfg, from, itemName)
	if err != nil {
		return nil, err
	}
	return FilterMap(found, nil, func(i int, e *T) string { return fieldStr(e, "Entity") }), nil
}

// from: existing, text, html
func Del[T any](cfg ItemConfig, from DbColType, itemName string) (int, error) {

	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, itemName)
	lk.Log("sFilter %v", sFilter)

	COL := cfg.DbColVal(from)
	if len(COL) == 0 {
		return 0, fmt.Errorf("from ItemDbCol can only be [existing, text, html]")
	}

	mh.UseDbCol(DATABASE, COL)
	if COL == cfg.DbColHtml {
		sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)>?%v<?"}}`, itemName)
	}
	nDeleted, _, err := mh.DeleteOne[T](strings.NewReader(sFilter)) // MUST re-create reader
	if err != nil {
		return 0, err
	}
	return nDeleted, nil
}

// from: existing, text, html
func Clr[T any](cfg ItemConfig, from DbColType) (int, error) {
	names, err := ListMany[T](cfg, from, "")
	if err != nil {
		return 0, err
	}
	cnt := 0
	for _, name := range names {
		n, err := Del[T](cfg, from, name)
		if err != nil {
			return 0, err
		}
		cnt += n
	}
	return cnt, nil
}

//////////////////////////////////////////////////////////////

func ColEntities(ColName string) ([]string, error) {

	mh.UseDbCol(DATABASE, "colentities") // fixed collection name

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

func EntClasses(EntName string) ([]string, []string, error) {

	mh.UseDbCol(DATABASE, "class") // fixed collection name

	whole, err := mh.FindOne[map[string]any](nil) // only one
	if err != nil {
		return nil, nil, err
	}
	rt, ok := (*whole)[EntName] // primitive.M
	if !ok {
		return []string{EntName}, []string{}, nil
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

func FullTextSearch(aim string, insensitive bool) ([]string, []string, error) {

	mh.UseDbCol(DATABASE, "pathval") // fixed collection name

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

////////////////////////////////////////////////////////////////////////////////

// to: submit approve
func RecordAction(user string, to DbColType, name, kind string) error {

	mh.UseDbCol(DATABASE, CfgAction.DbColVal(to))

	record, err := mh.FindOneAt[ActionRecord]("User", user)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	if record == nil {
		r := ActionRecord{
			User:   user,
			Action: string(to),
			Did: []DidItem{{
				Name:      name,
				Kind:      kind,
				Timestamp: time.Now(),
			}},
		}
		_, _, err := mh.Upsert(bytes.NewReader(r.Marshal()), "User", user)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return err
		}

	} else {

		record.Did = append(record.Did, DidItem{Name: name, Kind: kind, Timestamp: time.Now()})
		_, _, err := mh.Upsert(bytes.NewReader(record.Marshal()), "User", user)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return err
		}
	}
	return nil
}
