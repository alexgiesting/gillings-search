package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	db_connect()

	wd, _ := os.Getwd()
	log.Printf("working directory: %s\n", wd)

	PORT := ":8080"
	log.Printf("Running server on %s\n", PORT)
	// http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s %s", r.Method, r.URL.Path)))
	})
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func db_connect() {
	DB_HOST := os.Getenv("MONGODB_SERVICE_HOST")
	DB_PORT := os.Getenv("MONGODB_SERVICE_PORT")
	DB_USER := os.Getenv("MONGODB_USER")
	DB_PASSWORD := os.Getenv("MONGODB_PASSWORD")
	DB_NAME := os.Getenv("MONGODB_DATABASE")
	DB_URI := fmt.Sprintf("mongodb://%s:%s@%s:%s/", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		// TODO
		log.Fatal("@@@ failed to connect to MongoDB\n")
	}
	db := client.Database(DB_NAME)
	log.Printf("db: %v\n", db)
	client.Disconnect(context.TODO())
}
