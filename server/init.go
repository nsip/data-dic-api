package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	mh "github.com/digisan/db-helper/mongo"
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
	case <-time.After(3 * time.Second):
		fmt.Println("timed out for process")
	case x := <-ch:
		lk.FailOnErr("%v", x.err)

		// ingest existing entity json files
		ingestExistingEntities("./data/renamed", dbName, "entities")

		// ingestClassLinkage
		ingestClassLinkage("../data/out/class-link.json", dbName, "class")
	}
}

func ingestExistingEntities(dpath, db, col string) error {

	des, err := os.ReadDir(dpath)
	if err != nil {
		return err
	}

	mh.UseDbCol(db, col)

	nEntity := 0
	for _, de := range des {
		fpath := filepath.Join(dpath, de.Name())

		if !strings.HasSuffix(fpath, ".json") {
			continue
		}

		data, err := os.ReadFile(fpath)
		if err != nil {
			return err
		}
		entity := gjson.Get(string(data), "Entity").String()
		if len(entity) == 0 {
			return fmt.Errorf("entity value is empty, invalid file@ %v", fpath)
		}

		file, err := os.Open(fpath)
		if err != nil {
			return err
		}

		result, _, err := mh.Upsert(file, "Entity", entity)
		if err != nil {
			return err
		}
		lk.Log("ingesting... %v", result)

		if err := file.Close(); err != nil {
			return err
		}

		nEntity++
	}

	fmt.Printf("all %d entities have been ingested or replaced\n", nEntity)

	return nil
}

func ingestClassLinkage(fpath, db, col string) error {
	panic("TODO:")
	return nil
}
