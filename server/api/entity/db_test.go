package entity

import (
	"fmt"
	"testing"

	"github.com/nsip/data-dic-api/server/api/db"
)

func TestPingDB(t *testing.T) {

	userCollection := db.GetDbCol("testing", "users")

	rID, err := db.Insert(userCollection, `{
		"fullName": "User 1",
		"age": 1
	}`)

	if err != nil {
		panic(err)
	}
	fmt.Println(rID)

	////////////////////////////////////////////////////////

	rIDs, err := db.Insert(userCollection, `[
		{
			"fullName": "User 5",
			"age": 55
		},
		{
			"fullName": "User 6",
			"age": 66
		},
		{
			"fullName": "User 7",
			"age": 77
		}
	]`)

	if err != nil {
		panic(err)
	}
	fmt.Println(rIDs)
}
