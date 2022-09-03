package db

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

type Person struct {
	FullName string
	Age      int
	Class    struct {
		Name    string
		Teacher string
	}
}

func (p Person) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintln("FullName:", p.FullName))
	sb.WriteString(fmt.Sprintln("Age:", p.Age))
	sb.WriteString(fmt.Sprintln("Class.Name:", p.Class.Name))
	sb.WriteString(fmt.Sprintln("Class.Teacher:", p.Class.Teacher))
	return sb.String()
}

////////////////////////////////////////////////////////////////////////////////////////////////

func TestInsert(t *testing.T) {
	UseDbCol("testing", "users")

	////////////////////////////////////////////////////////

	r, err := os.Open("./s1.json")
	if err != nil {
		panic(err)
	}

	rID, err := Insert(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(rID)

	////////////////////////////////////////////////////////

	r, err = os.Open("./s2.json")
	if err != nil {
		panic(err)
	}

	rIDs, err := Insert(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(rIDs)
}

func TestFind(t *testing.T) {
	UseDbCol("testing", "users")

	////////////////////////////////////////////////////////

	// retrieve single and multiple documents with a specified filter using FindOne() and Find()
	// create a search filer
	// filter := bson.D{
	// 	{
	// 		"$and",
	// 		bson.A{
	// 			bson.D{
	// 				{
	// 					"age",
	// 					bson.D{{"$gt", 25}},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	rt, err := Find[Person](strings.NewReader(`{
		"$and": [
			{
				"age": {
					"$gt": 60
				}
			}
		]
	}`))

	// rt, err := Find[Person](nil)

	if err != nil {
		panic(err)
	}

	// fmt.Println(rt)

	for _, p := range rt {
		fmt.Println()
		fmt.Print(p)
	}
}

func TestUpdate(t *testing.T) {
	UseDbCol("testing", "users")

	rt, err := Update(
		strings.NewReader(`{
			"$and": [
				{
					"age": {
						"$gt": 60
					}
				}
			]
		}`),
		// nil,
		strings.NewReader(`{
			"$set": {
				"fullName": "User Modified"
			},
			"$inc": {
				"age": 1
			}
		}`),
		false,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(rt)
}

func TestDelete(t *testing.T) {
	UseDbCol("testing", "users")

	rt, err := Delete(
		strings.NewReader(`{
			"age": {
				"$lt": 50
			}				
		}`),
		false,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(rt)
}
