package db

import (
	"fmt"
	"testing"
)

func TestDel(t *testing.T) {
	fmt.Println(Del[any](CfgEnt, "test 3"))
	fmt.Println(Del[any](CfgCol, "test 3"))
}

func TestColEntities(t *testing.T) {
	entities, err := ColEntities(CfgCol, "NAPLAN Student Registration")
	fmt.Println(err)
	for _, en := range entities {
		fmt.Println(en)
	}
}

func TestEntClasses(t *testing.T) {
	derived, children, err := EntClasses(CfgEnt, "Staff")
	fmt.Println(err)
	fmt.Println("-----")
	for _, en := range derived {
		fmt.Println(en)
	}
	fmt.Println("-----")
	for _, en := range children {
		fmt.Println(en)
	}
}

func TestFullTextSearchh(t *testing.T) {
	ents, cols, err := FullTextSearch(CfgEnt, "http", false)
	fmt.Println(err)
	fmt.Println("-----")
	for _, ent := range ents {
		fmt.Println(ent)
	}
	fmt.Println("-----")
	for _, col := range cols {
		fmt.Println(col)
	}
}
