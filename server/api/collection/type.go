package collection

import (
	"fmt"
	"strings"
)

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
