package query

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/alexgiesting/gillings-search/database"
	"github.com/alexgiesting/gillings-search/paths"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type QueryHandler struct {
	db *mongo.Database
}

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// path := r.URL.Path[len(paths.PATH_QUERY):] // TODO
	query := r.URL.Query()

	filter := make(map[string]interface{})
	for k, v := range query {
		switch k {
		case "keyword":
			// TODO
		case "faculty":
			filter["authors.surname"] = v[0]
		case "dept":
			// TODO
		case "theme":
			// TODO
		}
	}
	bsonFilter, err := bson.Marshal(filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", filter) // TODO

	cursor, err := handler.db.Collection(database.CITATIONS).Find(context.TODO(), bsonFilter)
	if err != nil {
		log.Fatal(err)
	}
	var results []database.Citation
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Fatal(err)
	}
}

func Main() {
	client, db := database.Connect()
	defer client.Disconnect(context.TODO())

	serveMux := http.NewServeMux()
	serveMux.Handle(paths.PATH_QUERY, &QueryHandler{db})

	PORT := os.Getenv(paths.ENV_QUERY_PORT)
	log.Printf("Running server on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, serveMux))
}
