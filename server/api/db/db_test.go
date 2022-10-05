package db

import (
	"fmt"
	"testing"
)

func TestDel(t *testing.T) {
	fmt.Println(Del[any](CfgEntity, Existing, "test 3"))
	fmt.Println(Del[any](CfgCollection, Existing, "test 3"))
}

func TestColEntities(t *testing.T) {
	entities, err := GetColEntities("NAPLAN Student Registration")
	fmt.Println(err)
	for _, en := range entities {
		fmt.Println(en)
	}
}

func TestEntClasses(t *testing.T) {
	derived, children, err := GetEntClasses("Staff")
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
	ents, cols, err := FullTextSearch("http", false)
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
