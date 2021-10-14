package query

import (
	"context"
	"fmt"
	"html"
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
	path := r.URL.Path[len(paths.PATH_QUERY):] // TODO
	query := r.URL.Query()
	log.Printf("@@@ endpoint <%s%s>\n", path, r.URL.RawQuery)

	filter := make(map[string]interface{})
	for k, v := range query {
		switch k {
		case "keyword":
			// TODO
		case "faculty":
			if filter["authors"] == nil {
				filter["authors"] = make([]bson.M, 0, 1)
			}
			filter["authors"] = append(filter["authors"].([]bson.M), bson.M{"$elemMatch": bson.M{"surname": v}})
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
	log.Printf("%v", filter)

	cursor, err := handler.db.Collection(database.CITATIONS).Find(context.TODO(), bsonFilter)
	if err != nil {
		log.Fatal(err)
	}
	var results []database.Citation
	cursor.All(context.TODO(), &results)

	resultsString := html.EscapeString(fmt.Sprintf("%v", results))
	message := fmt.Sprintf("<!DOCTYPE html><html><body><pre>%v</pre></body></html>", resultsString)
	w.Write([]byte(message))
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
