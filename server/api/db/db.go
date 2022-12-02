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

// all Item in db must have 'Entity' field !!!

// from: existing, text, html; [itemName] is Item 'Entity' value
func One[T any](cfg ItemConfig, from DbColType, itemName string, fuzzy bool) (*T, error) {
	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, itemName)
	if fuzzy {
		sFilter = fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, itemName)
	}
	lk.Log("sFilter %v", sFilter)

	mh.UseDbCol(DATABASE, string(cfg.DbColName(from)))

	found, err := mh.FindOne[T](strings.NewReader(sFilter))
	if err != nil {
		return nil, err
	}
	return found, nil
}

// from: existing, text, html; [filterName] is Item 'Entity' value
func Many[T any](cfg ItemConfig, from DbColType, filterName string) ([]*T, error) {
	var rFilter io.Reader = nil // NOT using "*strings.Reader = nil" as nil interface
	if len(filterName) > 0 {
		sFilter := fmt.Sprintf(`{"Entity": {"$regex": "(?i)%v"}}`, filterName)
		lk.Log("sFilter %v", sFilter)
		rFilter = strings.NewReader(sFilter)
	}

	mh.UseDbCol(DATABASE, string(cfg.DbColName(from)))

	found, err := mh.Find[T](rFilter)
	if err != nil {
		return nil, err
	}
	return found, nil
}

// from: existing, text, html; return list of Item 'Entity' value
func ListMany[T any](cfg ItemConfig, from DbColType, filterName string) ([]string, error) {
	found, err := Many[T](cfg, from, filterName)
	if err != nil {
		return nil, err
	}
	return FilterMap(found, nil, func(i int, e *T) string { return fieldStr(e, "Entity") }), nil
}

// from: existing, text, html, [itemName] is Item 'Entity' value
func Del[T any](cfg ItemConfig, from DbColType, itemName string) (int, error) {

	sFilter := fmt.Sprintf(`{"Entity": "%v"}`, itemName)
	lk.Log("sFilter %v", sFilter)

	COL := cfg.DbColName(from)
	mh.UseDbCol(DATABASE, string(COL))

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

func GetColEntities(ColName string) ([]string, error) {

	mh.UseDbCol(DATABASE, string(ColEntities)) // fixed collection name

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

func GetEntClasses(EntName string) ([]string, []string, error) {

	mh.UseDbCol(DATABASE, string(Class)) // fixed collection name

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

	mh.UseDbCol(DATABASE, string(PathVal)) // fixed collection name

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

// from: existing, text, html
func Exists(from DbColType, name string) (bool, error) {
	for kind, cfg := range CfgGrp {
		var (
			result any = nil
			err    error
		)
		switch kind {
		case "entity":
			result, err = One[EntType](cfg, from, name, false)
		case "collection":
			result, err = One[ColType](cfg, from, name, false)
		}
		if err != nil {
			lk.WarnOnErr("%v", err)
			return false, err
		}
		if !IsNil(result) {
			return true, nil
		}
	}
	return false, nil
}

////////////////////////////////////////////////////////////////////////////////

// [from] DbColType can only be [submit, approve, subscribe]
func ActionExists(user string, from DbColType) (bool, error) {

	mh.UseDbCol(DATABASE, string(CfgAction.DbColName(from)))

	record, err := mh.FindOneAt[ActionRecord]("User", user)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	if record == nil {
		return false, nil
	}
	return true, nil
}

func ActionRecordExists(user string, from DbColType, name string) (bool, error) {

	ok, err := ActionExists(user, from)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	if !ok {
		return false, nil
	}

	record, _ := mh.FindOneAt[ActionRecord]("User", user)
	for _, item := range record.Did {
		if item.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// to: [submit approve subscribe]
func RecordAction(user string, to DbColType, name, kind string) (bool, error) {

	// check item exists under this action
	ok, err := ActionRecordExists(user, to, name)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	// already exists, do nothing
	if ok {
		return false, nil
	}

	// check action exists
	ok, err = ActionExists(user, to)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	if ok { // action exists, append item

		record, _ := mh.FindOneAt[ActionRecord]("User", user)
		record.Did = append(record.Did, DidItem{Name: name, Kind: kind, Timestamp: time.Now()})
		if _, _, err = mh.Upsert(bytes.NewReader(record.Marshal()), "User", user); err != nil {
			lk.WarnOnErr("%v", err)
			return false, err
		}

	} else { // action doesn't exist, create new

		r := ActionRecord{
			User:   user,
			Action: string(to),
			Did: []DidItem{{
				Name:      name,
				Kind:      kind,
				Timestamp: time.Now(),
			}},
		}
		if _, _, err := mh.Upsert(bytes.NewReader(r.Marshal()), "User", user); err != nil {
			lk.WarnOnErr("%v", err)
			return false, err
		}
	}
	return true, nil
}

// to: [submit approve subscribe]
func RemoveAction(user string, from DbColType, name string) (bool, error) {

	// check action exists
	ok, err := ActionExists(user, from)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	if !ok { // action doesn't exist, do nothing
		return false, nil
	}

	// check item exists under this action
	ok, err = ActionRecordExists(user, from, name)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	if !ok { // action item doesn't exist, do nothing
		return false, nil
	}

	record, _ := mh.FindOneAt[ActionRecord]("User", user)
	record.Did = Filter(record.Did, func(i int, e DidItem) bool { return e.Name != name })
	if _, _, err = mh.Upsert(bytes.NewReader(record.Marshal()), "User", user); err != nil {
		lk.WarnOnErr("%v", err)
		return false, err
	}
	return true, nil
}

func ListActionRecord(user string, from DbColType) ([]string, error) {

	// check action exists
	ok, err := ActionExists(user, from)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return nil, err
	}
	if !ok { // action doesn't exist, return empty
		return []string{}, nil
	}

	record, err := mh.FindOneAt[ActionRecord]("User", user)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return []string{}, err
	}
	return FilterMap(record.Did, nil, func(i int, e DidItem) string { return e.Name }), nil
}
