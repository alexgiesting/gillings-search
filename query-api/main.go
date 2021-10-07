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

type QueryHandler struct {
	db *mongo.Database
}

func main() {
	client, db := db_connect()
	defer client.Disconnect(context.TODO())
	http.Handle("/", &QueryHandler{db})

	PORT := ":8080"
	log.Printf("Running server on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func db_connect() (*mongo.Client, *mongo.Database) {
	DB_HOST := os.Getenv("MONGODB_SERVICE_HOST")
	DB_PORT := os.Getenv("MONGODB_SERVICE_PORT")
	DB_USER := os.Getenv("MONGODB_USER")
	DB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	DB_NAME := os.Getenv("MONGODB_DATABASE")
	DB_URI := fmt.Sprintf("mongodb://%s:%s@%s:%s/", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatal("@@@ failed to connect to MongoDB\n") // TODO
	}
	db := client.Database(DB_NAME)
	return client, db
}

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path[1:]
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
