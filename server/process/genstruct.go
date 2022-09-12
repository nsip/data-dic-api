package process

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

func GenEntityPathVal(fpaths ...string) (map[string]string, error) {
	m := make(map[string]string)
	for _, fpath := range fpaths {
		if strs.HasAnySuffix(fpath, "class-link.json", "collection-entities.json") {
			continue
		}
		data, err := os.ReadFile(fpath)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return nil, err
		}

		mPathVal, err := jt.Flatten(data)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return nil, err
		}

		key, ok := mPathVal["Entity"]
		if !ok {
			lk.WarnOnErr("%v @ "+fpath, errors.New("entity missing"))
			return nil, fmt.Errorf("%v @ "+fpath, "entity missing")
		}

		// make json
		js := "{"
		for path, val := range mPathVal {
			path = strings.ReplaceAll(path, `.`, `[dot]`)
			val = strings.ReplaceAll(val.(string), `"`, `\"`)
			js += fmt.Sprintf(`"%s": "%s",`, path, val)
		}
		js = strings.TrimSuffix(js, ",") + "}"
		if !jt.IsValidStr(js) {
			lk.WarnOnErr("%v @"+fpath, errors.New("invalid path-value json"))
			return nil, fmt.Errorf("%v @"+fpath, "invalid path-value json")
		}

		m[key.(string)] = js
	}
	return m, nil
}

func DumpPathValue(idir, odname string) error {
	osdir := filepath.Join(idir, odname)
	gio.MustCreateDir(osdir)
	fpaths, _, err := fd.WalkFileDir(idir, false)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	mEntPathVal, err := GenEntityPathVal(fpaths...)
	if err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	for entity, js := range mEntPathVal {
		err := os.WriteFile(filepath.Join(osdir, entity+".json"), []byte(js), os.ModePerm)
		if err != nil {
			lk.WarnOnErr("%v", err)
			return err
		}
	}
	return nil
}
