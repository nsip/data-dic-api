package process

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func DumpCollection(dir, ofname, idfield string, idvalue any) error {
	mColEntities := make(map[string][]string)
	fis, err := os.ReadDir(dir)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}
	for _, fi := range fis {
		if fname := fi.Name(); strings.HasSuffix(fname, ".json") {
			bytesJS, err := os.ReadFile(filepath.Join(dir, fname))
			if err != nil {
				lk.WarnOnErr("%v", err)
				return err
			}
			js := string(bytesJS)
			if r := gjson.Get(js, "Collections"); r.IsArray() {
				for i := 0; i < len(r.Array()); i++ {
					colName := gjson.Get(js, fmt.Sprintf("Collections.%d.Name", i)).String()
					mColEntities[colName] = append(mColEntities[colName], strings.TrimSuffix(fname, ".json"))
				}
			}
		}
	}
	bytesJS, err := json.Marshal(mColEntities)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}
	js, err := sjson.Set(string(bytesJS), idfield, idvalue)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}
	err = os.WriteFile(filepath.Join(dir, ofname), []byte(js), os.ModePerm)
	lk.WarnOnErr("%v", err)
	return err
}
