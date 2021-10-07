package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	PORT = ":8080"
	PATH = "/query"
)

type QueryHandler struct {
	db *mongo.Database
}

func main() {
	client, db := db_connect()
	defer client.Disconnect(context.TODO())
	db_init(db)
	http.Handle("/", &QueryHandler{db})

	log.Printf("Running server on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func db_connect() (*mongo.Client, *mongo.Database) {
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
	return client, db
}

func db_init(db *mongo.Database) {
	b_str := func(description string) bson.M { return bson.M{"bsonType": "string", "description": description} }
	b_int := func(description string) bson.M { return bson.M{"bsonType": "int", "description": description} }

	add_collection(db, "faculty", bson.M{
		"name": b_str("faculty name"),
		"sid":  b_int("Scopus ID"),
	})

	add_collection(db, "publications", bson.M{
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

func add_collection(db *mongo.Database, name string, schema bson.M) {
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

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path[len(PATH)+1:]
	collections, err := handler.db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Print(err) // TODO
	}
	found := false
	for _, name := range collections {
		if query == name {
			found = true
			break
		}
	}

	var message string
	if found {
		message = fmt.Sprintf("`%s` found in collections", query)
	} else {
		message = fmt.Sprintf("`%s` not found in collections", query)
	}
	message = fmt.Sprintf("%s %s\n\n%s: {%s}", r.Method, r.URL.Path, message, strings.Join(collections, ","))
	message = fmt.Sprintf("<!DOCTYPE html><html><body><pre>%s</pre></body></html>", message)
	w.Write([]byte(message))
}
