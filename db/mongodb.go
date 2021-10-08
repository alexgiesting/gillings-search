package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	mongo.Database
}

func Connect() (*mongo.Client, *Database) {
	DB_HOST := os.Getenv("MONGODB_SERVICE_HOST")
	DB_PORT := os.Getenv("MONGODB_SERVICE_PORT")
	// DB_USER := os.Getenv("MONGODB_USER")
	// DB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	DB_ADMIN_PASSWORD := os.Getenv("MONGODB_ADMIN_PASSWORD")
	DB_NAME := os.Getenv("MONGODB_DATABASE")
	DB_URI := fmt.Sprintf("mongodb://%s:%s@%s:%s/", "admin", DB_ADMIN_PASSWORD, DB_HOST, DB_PORT)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err) // TODO
	}
	db := client.Database(DB_NAME)
	return client, &Database{*db}
}

func (db *Database) Init() {
	b_str := func(description string) bson.M { return bson.M{"bsonType": "string", "description": description} }
	b_int := func(description string) bson.M { return bson.M{"bsonType": "int", "description": description} }

	db.AddCollection("faculty", bson.M{
		"name": b_str("faculty name"),
		"sid":  b_int("Scopus ID"),
	})

	db.AddCollection("publications", bson.M{
		"authors": bson.M{
			"description": "author names",
			"bsonType":    "array",
			"items":       b_str("author name"),
		},
		"title":  b_str("title"),
		"year":   b_int("year"),
		"source": b_str("source title"),
		"volume": b_str("source volume"),
		"issue":  b_str("source issue"),
		"number": b_str("number"),
		"doi":    b_str("DOI"),
		"eid":    b_str("Scopus ID"),
	})
}

func (db *Database) AddCollection(name string, schema bson.M) {
	log.Printf("Checking for collection `%s`...\n", name)
	err := db.CreateCollection(context.TODO(), name, options.CreateCollection().SetValidator(bson.M{
		"$jsonSchema": bson.M{"bsonType": "object", "properties": schema},
	}))
	if err == nil {
		log.Printf("Added `%s` to the database\n", name)
	} else {
		log.Printf("Failed to add `%s`: %v\n", name, err)
	}
}
