package main

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

func init() {

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
		fmt.Println("timed out for process")
	case x := <-ch:
		lk.FailOnErr("%v", x.err)

		// ingest existing entities json files
		lk.FailOnErr("%v", ingestFromDir(dbName, "entities", "./data/out", "Entity", "class-link.json", "collection-entities.json"))

		// ingest Entity ClassLinkage
		lk.FailOnErr("%v", ingestFromFile(dbName, "class", "./data/out/class-link.json", "ver"))

		// ingest Entities PathVal
		lk.FailOnErr("%v", ingestFromDir(dbName, "pathval", "./data/out/path_val", "Entity"))

		//////////////////////////////

		// ingest Collections
		lk.FailOnErr("%v", ingestFromDir(dbName, "collections", "./data/out/collections", "Entity", "class-link.json", "collection-entities.json"))

		// ingest Collections PathVal
		lk.FailOnErr("%v", ingestFromDir(dbName, "pathval", "./data/out/collections/path_val", "Entity"))

		// ingest Collection-Entities
		lk.FailOnErr("%v", ingestFromFile(dbName, "colentities", "./data/out/collection-entities.json", "ver"))
	}
}

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

		result, _, err := mh.Upsert(file, idfield, id)
		if err != nil {
			return err
		}
		lk.Log("ingesting... %v", result)

		if err := file.Close(); err != nil {
			return err
		}

		nFile++
	}

	lk.Log("all [%d] files have been ingested or updated\n", nFile)
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

	lk.Log("[%v] has been updated", col)
	return nil
}
