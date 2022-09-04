package collection

import (
	"strings"
)

type CollectionType struct {
}

func (c CollectionType) String() string {
	sb := strings.Builder{}
	// sb.WriteString(fmt.Sprintln("Entity:", e.Entity))
	return sb.String()
}
