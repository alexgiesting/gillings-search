package query

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/alexgiesting/gillings-search/go/database"
	"github.com/alexgiesting/gillings-search/go/paths"
	"go.mongodb.org/mongo-driver/bson"
)

type QueryHandler struct {
	db *database.Connection
}

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// path := r.URL.Path[len(paths.PATH_QUERY):] // TODO
	query := r.URL.Query()
	var querySearch map[string]interface{}
	json.Unmarshal([]byte(query.Get("q")), &querySearch)

	// search := handler.db.Citations.Filter(makeSearch(querySearch))
	var results []database.Citation
	err := handler.db.Citations.Decode(&results) // search.Decode(&results)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(results)
	// err = json.NewEncoder(w).Encode(struct{ Results []database.Citation }{Results: results}),
	if err != nil {
		log.Fatal(err)
	}
}

func d(key string, value interface{}) bson.D {
	return bson.D{{Key: key, Value: value}}
}

func makeSearch(document map[string]interface{}) (string, []bson.D) {
	// TODO limit/paginate results
	// TODO validate fields
	filter := []bson.D{}
	for k, v := range document {
		switch k {
		case "or":
			orFilter := []bson.D{}
			for _, w := range v.([]interface{}) {
				orDocument := w.(map[string]interface{})
				orFilter = append(orFilter, d(makeSearch(orDocument)))
			}
			filter = append(filter, d("$or", orFilter))
		case "keyword":
			filter = append(filter, match(v, "keywords", "title", "abstract"))
		case "faculty":
			filter = append(filter, match(v, "authors.name"))
		// case "dept": // TODO
		// case "theme": // TODO
		default:
			log.Printf("Unrecognized parameter `%s` with value: %v", k, v)
		}
	}
	return "$and", filter
}

func match(v interface{}, fields ...string) bson.D {
	words := v.([]interface{})
	if len(words) == 1 {
		return matchElement(words[0].(string), fields)
	} else {
		filter := make([]bson.D, len(words))
		for i, word := range words {
			filter[i] = matchElement(word.(string), fields)
		}
		return d("$and", filter)
	}
}

func matchElement(word string, fields []string) bson.D {
	// TODO figure out how text indices work :(
	matchWord := d("$text", d("$search", word))
	if len(fields) == 1 {
		return d(fields[0], matchWord)
	} else {
		alternatives := make([]bson.D, len(fields))
		for j, field := range fields {
			alternatives[j] = d(field, matchWord)
		}
		return d("$or", alternatives)
	}
}

func Main() {
	db := database.Connect()
	defer db.Disconnect(context.TODO())

	serveMux := http.NewServeMux()
	serveMux.Handle(paths.PATH_QUERY, &QueryHandler{db})

	PORT := os.Getenv(paths.ENV_QUERY_PORT)
	log.Printf("Running server on %s", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, serveMux))
}
