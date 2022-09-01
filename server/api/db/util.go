package db

import (
	"fmt"
	"strings"

	lk "github.com/digisan/logkit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func GetDbCol(dbName, colName string) *mongo.Collection {
	return Client.Database(dbName).Collection(colName)
}

// user := bson.D{
// 	{"fullName", "User 1"},
// 	{"age", 32},
// }

func Insert(col *mongo.Collection, json string) (any, error) {

	json = strings.TrimSpace(json)
	if len(json) == 0 {
		return nil, nil
	}

	isArray := false
	if json[0] == '[' {
		isArray = true
	}

	if isArray {

		var docs []any
		err := bson.UnmarshalExtJSON([]byte(json), true, &docs)
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
		err := bson.UnmarshalExtJSON([]byte(json), true, &doc)
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
