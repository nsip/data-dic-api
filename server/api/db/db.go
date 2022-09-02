package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ref https://blog.logrocket.com/how-to-use-mongodb-with-go/

var (
	col *mongo.Collection
)

func getMongoClient(ip string, port int) *mongo.Client {
	if len(ip) == 0 {
		ip = "localhost"
	}
	if port == 0 {
		port = 27017
	}
	uri := fmt.Sprintf("mongodb://%s:%d", ip, port)
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))
	lk.FailOnErr("Connect error: %v", err)
	lk.FailOnErr("Ping error: %v", client.Ping(Ctx, readpref.Primary()))
	return client
}

func UseDbCol(dbName, colName string) {
	col = Client.Database(dbName).Collection(colName)
}

// return json string, is array type, error
func reader4json(r io.Reader) ([]byte, bool, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, false, err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return []byte{}, false, nil
	}
	if !jt.IsValid(data) {
		return nil, false, fmt.Errorf("invalid JSON")
	}
	return data, data[0] == '[', nil
}

func Insert(rData io.Reader) (any, error) {

	lk.FailOnErrWhen(col == nil, "%v", fmt.Errorf("collection is nil, use 'UseDbCol' to init one"))

	if rData == nil {
		return 0, nil
	}

	dataJSON, isArray, err := reader4json(rData)
	if err != nil {
		return nil, err
	}

	if isArray {

		var docs []any
		err := bson.UnmarshalExtJSON(dataJSON, true, &docs)
		if err != nil {
			return nil, err
		}
		result, err := col.InsertMany(Ctx, docs)
		if err != nil {
			return nil, err
		}
		return result.InsertedIDs, nil

	} else {

		var doc any
		err := bson.UnmarshalExtJSON(dataJSON, true, &doc)
		if err != nil {
			return nil, err
		}
		result, err := col.InsertOne(Ctx, doc)
		if err != nil {
			return nil, err
		}
		return result.InsertedID, nil
	}
}

func Find[T any](rFilter io.Reader) (rt []*T, err error) {

	lk.FailOnErrWhen(col == nil, "%v", fmt.Errorf("collection is nil, use 'UseDbCol' to init one"))

	var filter any

	if rFilter != nil {
		filterJSON, _, err := reader4json(rFilter)
		if err != nil {
			return nil, err
		}
		if err := bson.UnmarshalExtJSON(filterJSON, true, &filter); err != nil {
			return nil, err
		}
	} else {
		filter = bson.D{}
	}

	cursor, err := col.Find(Ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err = cursor.All(Ctx, &results); err != nil {
		return nil, err
	}

	for _, r := range results {
		data, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}

		one := new(T)
		err = json.Unmarshal(data, one)
		if err != nil {
			return nil, err
		}

		rt = append(rt, one)
	}

	return rt, nil
}

func FindOne[T any](rFilter io.Reader) (*T, error) {

	lk.FailOnErrWhen(col == nil, "%v", fmt.Errorf("collection is nil, use 'UseDbCol' to init one"))

	var filter any

	if rFilter != nil {
		filterJSON, _, err := reader4json(rFilter)
		if err != nil {
			return nil, err
		}
		if err := bson.UnmarshalExtJSON(filterJSON, true, &filter); err != nil {
			return nil, err
		}
	} else {
		filter = bson.D{}
	}

	one := new(T)
	if err := col.FindOne(Ctx, filter).Decode(one); err != nil {
		return nil, err
	}

	return one, nil
}

///////////////////////////////////////////////

// func Update(r io.Reader) (any, error) {
// 	lk.FailOnErrWhen(col == nil, "%v", fmt.Errorf("collection is nil, use 'UseDbCol' to init one"))

// 	json, isArray, err := payload4json(r)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if isArray {

// 		var docs []any
// 		err := bson.UnmarshalExtJSON([]byte(json), true, &docs)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result, err := col.UpdateMany(Ctx) // InsertMany(Ctx, docs)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return result.UpsertedCount, nil

// 	} else {

// 		var doc any
// 		err := bson.UnmarshalExtJSON([]byte(json), true, &doc)
// 		if err != nil {
// 			return nil, err
// 		}
// 		result, err := col.UpdateOne(Ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return result.UpsertedCount, nil
// 	}
// }
