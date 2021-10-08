package query

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alexgiesting/gillings-search/db"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	PORT = ":8080"
	PATH = "/query"
)

type QueryHandler struct {
	db *db.Database
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

func Main() {
	client, db := db.Connect()
	defer client.Disconnect(context.TODO())
	db.Init()
	http.Handle("/", &QueryHandler{db})

	log.Printf("Running server on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
