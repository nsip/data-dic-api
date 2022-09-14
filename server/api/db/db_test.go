package db

import (
	"fmt"
	"testing"
)

func TestDel(t *testing.T) {
	fmt.Println(Del[any](CfgEnt, "test 3"))
	fmt.Println(Del[any](CfgCol, "test 3"))
}
