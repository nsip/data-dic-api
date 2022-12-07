package db

import (
	"encoding/json"
	"time"

	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

// all Item must have 'Metadata.Type' field !!!

// same as ingested json schema
type EntType struct {
	Entity     string // *
	OtherNames []string
	Definition string
	SIF        []struct {
		XPath      []string
		Definition string
		Commentary string
		Datestamp  string
	}
	OtherStandards []struct {
		Standard   string
		Link       []string
		Path       []string
		Definition string
		Commentary string
	}
	LegalDefinitions []struct {
		LegislationName string
		Citation        string
		Link            string
		Definition      string
		Commentary      string
		Datestamp       string
	}
	Collections []struct {
		Name                   string
		Description            string
		Standard               string
		Elements               []string
		BusinessRules          []string
		DefinitionModification string
	}
	Metadata struct {
		Identifier         string
		Type               string //*
		ExpectedAttributes []string
		Superclass         []string
		CrossrefEntities   []string
	}
}

// same as ingested json schema
type ColType struct {
	Entity     string // *
	Definition string
	URL        []string
	Metadata   struct {
		Identifier string
		Type       string // *
	}
	Entities []string // this field will be filled from other dbcol
}

func ItemKind(data []byte) string {
	var (
		ent = EntType{}
		col = ColType{}
	)
	switch {
	case json.Unmarshal(data, &ent) == nil:
		return "entity"
	case json.Unmarshal(data, &col) == nil:
		return "collection"
	default:
		mt := gjson.Get(string(data), "Metadata.Type").String()
		lk.Warn("%v", mt)
		return mt
	}
}

func Item[T any](data []byte) *T {
	var (
		ent any = &EntType{}
		col any = &ColType{}
	)
	switch {
	case json.Unmarshal(data, ent) == nil:
		return ent.(*T)
	case json.Unmarshal(data, col) == nil:
		return col.(*T)
	default:
		return nil
	}
}

/////////////////////////////////////////////////

type DidItem struct {
	Name      string    // "Entity" value
	Kind      string    // Entity or Collection
	Timestamp time.Time //
}

type ActionRecord struct {
	User   string    // uname
	Action string    // [submit, approve]
	Did    []DidItem // list of Did
}

func (r *ActionRecord) Marshal() []byte {
	data, err := json.Marshal(r)
	lk.FailP1OnErr("%v", err)
	return data
}
