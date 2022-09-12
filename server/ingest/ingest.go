package ingest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

const (
	dbName = "dictionaryTest"
)

func IngestViaCmd() error {

	type output struct {
		out []byte
		err error
	}

	ch := make(chan output)
	go func() {
		// invoke 'process'
		cmd := exec.Command("./process") // without arguments, only process from 'data/renamed'
		out, err := cmd.CombinedOutput() // if successful, './data/out' will be created
		ch <- output{out, err}
	}()

	select {
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timed out for ingestion")

	case x := <-ch:
		lk.WarnOnErr("exec.Command error: [%v]", x.err)
		if x.err != nil {
			return x.err
		}
		if err := ingestAll(); err != nil {
			return err
		}
	}

	return nil
}

func ingestAll() error {

	// ingest existing entities json files
	if err := ingestFromDir(dbName, "entities", "./data/out", "Entity", "class-link.json", "collection-entities.json"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Entity ClassLinkage
	if err := ingestFromFile(dbName, "class", "./data/out/class-link.json", "RefName"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Entities PathVal
	if err := ingestFromDir(dbName, "pathval", "./data/out/path_val", "Entity"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	//////////////////////////////

	// ingest Collections
	if err := ingestFromDir(dbName, "collections", "./data/out/collections", "Entity", "class-link.json", "collection-entities.json"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Collections PathVal
	if err := ingestFromDir(dbName, "pathval", "./data/out/collections/path_val", "Entity"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Collection-Entities
	if err := ingestFromFile(dbName, "colentities", "./data/out/collection-entities.json", "RefName"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ingestFromDir(db, col, dpath, idfield string, exclfiles ...string) error {

	des, err := os.ReadDir(dpath)
	if err != nil {
		return err
	}

	mh.UseDbCol(db, col)

	nFile := 0
	for _, de := range des {
		fpath := filepath.Join(dpath, de.Name())

		if !strings.HasSuffix(fpath, ".json") {
			continue
		}
		if In(de.Name(), exclfiles...) {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}

		id := gjson.Get(string(data), idfield).String()
		if len(id) == 0 {
			lk.Warn("ID(%v) value is empty, file@ %v, ignored", idfield, fpath)
			continue
		}

		file, err := os.Open(fpath)
		if err != nil {
			return err
		}

		_, _, err = mh.Upsert(file, idfield, id)
		if err != nil {
			return err
		}
		lk.Log("ingesting... %v", id)

		if err := file.Close(); err != nil {
			return err
		}

		nFile++
	}

	lk.Log("all [%d] files have been ingested or updated on [%v]", nFile, col)
	return nil
}

func ingestFromFile(db, col, fpath, idfield string) error {

	if !strings.HasSuffix(fpath, ".json") {
		return nil
	}

	data, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}

	id := gjson.Get(string(data), idfield).String()
	if len(id) == 0 {
		lk.Log("ID(%v) value is empty, file@ %v, to do insert", idfield, fpath)
	}

	file, err := os.Open(fpath)
	if err != nil {
		return err
	}

	mh.UseDbCol(db, col)

	_, _, err = mh.Upsert(file, idfield, id)
	if err != nil {
		return err
	}

	lk.Log("[%v] has been updated on [%v]", id, col)
	return nil
}
