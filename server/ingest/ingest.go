package ingest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	mh "github.com/digisan/db-helper/mongo"
	. "github.com/digisan/go-generics/v2"
	lk "github.com/digisan/logkit"
	"github.com/nsip/data-dic-api/server/api/db"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	dbName = "MyDictionary"
)

func IngestViaCmd(clrdb bool) error {

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
		if clrdb {
			if err := clearDb(dbName); err != nil {
				return err
			}
		}
		if err := ingestAll(dbName); err != nil {
			return err
		}
	}

	return nil
}

func clearDb(DB string) error {
	if err := mh.DropCol(DB, "pathval"); err != nil {
		return err
	}
	if err := mh.DropCol(DB, "class"); err != nil {
		return err
	}
	if err := mh.DropCol(DB, "colentities"); err != nil {
		return err
	}
	if err := mh.DropCol(DB, "collections"); err != nil {
		return err
	}
	if err := mh.DropCol(DB, "entities"); err != nil {
		return err
	}
	return nil
}

func ingestAll(DB string) error {

	// ingest existing entities json files
	if err := ingestFromDir(DB, "entities", "./data/out", "Entity", "class-link.json", "collection-entities.json"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Entity ClassLinkage
	if err := ingestFromFile(DB, "class", "./data/out/class-link.json", "RefName"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Entities PathVal
	if err := ingestFromDir(DB, "pathval", "./data/out/path_val", "Entity"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Collections
	if err := ingestFromDir(DB, "collections", "./data/out/collections", "Entity", "class-link.json", "collection-entities.json"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Collections PathVal
	if err := ingestFromDir(DB, "pathval", "./data/out/collections/path_val", "Entity"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// ingest Collection-Entities
	if err := ingestFromFile(DB, "colentities", "./data/out/collection-entities.json", "RefName"); err != nil {
		lk.WarnOnErr("%v", err)
		return err
	}

	// append 'colentities' to 'collections' 'Entities' field
	colEntitiesToCollections(DB)

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ingestFromDir(DB, Col, dpath, idfield string, exclfiles ...string) error {

	des, err := os.ReadDir(dpath)
	if err != nil {
		return err
	}

	mh.UseDbCol(DB, Col)

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

	lk.Log("all [%d] files have been ingested or updated on [%v]", nFile, Col)
	return nil
}

func ingestFromFile(DB, Col, fpath, idfield string) error {

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

	mh.UseDbCol(DB, Col)

	_, _, err = mh.Upsert(file, idfield, id)
	if err != nil {
		return err
	}

	lk.Log("[%v] has been updated on [%v]", id, Col)
	return nil
}

func colEntitiesToCollections(DB string) {

	mh.UseDbCol(DB, "colentities") // this db-collection only has one single doc

	colEntities, err := mh.FindOneAt[map[string]any]("RefName", "CollectionEntities")
	lk.FailOnErr("%v", err)

	delete(*colEntities, "_id")
	delete(*colEntities, "RefName") // "CollectionEntities" must have field "RefName"

	mh.UseDbCol(DB, "collections") // in this db-collection, each Collection-Item has its own doc

	collections, err := mh.Find[db.ColType](nil)
	lk.FailOnErr("%v", err)
	for _, col := range collections {

		for _, e := range (*colEntities)[col.Entity].(primitive.A) {
			col.Entities = append(col.Entities, e.(string))
		}

		if len(col.Entities) > 0 {
			data, err := json.Marshal(*col)
			lk.FailOnErr("%v", err)

			_, _, err = mh.ReplaceOneAt("Entity", col.Entity, bytes.NewReader(data))
			lk.FailOnErr("%v", err)
		}
	}
}
