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
)

type QueryHandler struct {
	db *database.Connection
}

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// path := r.URL.Path[len(paths.PATH_QUERY):] // TODO
	query := r.URL.Query()

	search := handler.db.Citations.Search()
	for k, v := range query {
		switch k {
		case "keyword":
			if len(v) == 1 && v[0] == "" { // TODO
				break
			}
			search.Filter("keywords", bson.M{"$in": v})
		case "faculty":
			if len(v) == 1 && v[0] == "" { // TODO
				break
			}
			search.Filter("authors.surname", bson.M{"$all": v})
		case "dept": // TODO
		case "theme": // TODO
		}
	}
	log.Printf("%v", search) // TODO

	var results []database.Citation
	err := search.Decode(&results)
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
	db := database.Connect()
	defer db.Disconnect(context.TODO())

	serveMux := http.NewServeMux()
	serveMux.Handle(paths.PATH_QUERY, &QueryHandler{db})

	PORT := os.Getenv(paths.ENV_QUERY_PORT)
	log.Printf("Running server on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, serveMux))
}
