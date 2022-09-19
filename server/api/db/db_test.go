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