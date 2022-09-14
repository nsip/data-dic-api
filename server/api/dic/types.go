package dic

import (
	"fmt"
	"strings"
)

type EntityType struct {
	Entity     string
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
		DefinitionModification string
	}
	Metadata struct {
		Identifier         string
		Type               string
		ExpectedAttributes []string
		Superclass         []string
		CrossrefEntities   []string
	}
}

func (e EntityType) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintln("Entity:", e.Entity))
	return sb.String()
}

type CollectionType struct {
	Entity     string
	Definition string
	URL        []string
	Metadata   struct {
		Identifier string
		Type       string
	}
}

func (c CollectionType) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintln("Entity:", c.Entity))
	return sb.String()
}
