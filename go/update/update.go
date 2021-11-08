package update

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/alexgiesting/gillings-search/go/database"
	"github.com/alexgiesting/gillings-search/go/paths"
)

func update(db *database.Connection, r Request) {
	switch r.path {
	case "pull":
		// TODO where should the date come from?
		startDate := "2021-01-01"
		pullCitations(db, startDate)
	case "push":
		pushCitations(db)
	case "load/faculty":
		db.Faculty.Drop(context.TODO())
		db.LoadFaculty(r.body)
	case "drop/faculty":
		db.Faculty.Drop(context.TODO())
	case "load/citations":
		if r.query.Get("drop") == "1" {
			db.Citations.Drop(context.TODO())
		}
		db.LoadCitations(r.body)
	case "drop/citations":
		db.Citations.Drop(context.TODO())
	case "load/themes":
		db.Themes.Drop(context.TODO())
		db.LoadThemes(r.body)
	case "drop/themes":
		db.Themes.Drop(context.TODO())
	default:
		log.Printf("Invalid request `%s` received by `update`", r.path)
	}
}

type QueryHandler struct {
	updateKey string
	request   chan Request
}

type Request struct {
	path  string
	query url.Values
	body  *bytes.Reader
}

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO maybe not form?
	if r.FormValue("key") != handler.updateKey {
		// TODO use userinfo instead?
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	path := strings.TrimRight(r.URL.Path[len(paths.PATH_UPDATE):], "/")
	var body *bytes.Reader
	if path[:len("load")] == "load" { // TODO I think this throws if the path is too short?
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Fatal(err)
		}
		bodyBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
		body = bytes.NewReader(bodyBytes)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "command received")

	handler.request <- Request{path, r.Form, body}
}

func Main() {
	db := database.Connect()
	defer db.Disconnect(context.TODO())

	updateKey, err := paths.LoadKey(paths.SECRET_UPDATE_KEY)
	if err != nil {
		log.Fatal(err)
	}
	request := make(chan Request)

	serveMux := http.NewServeMux()
	handler := QueryHandler{updateKey, request}
	serveMux.Handle(paths.PATH_UPDATE, &handler)
	PORT := os.Getenv(paths.ENV_UPDATE_PORT)
	log.Printf("Running server on %s", PORT)
	go func() { log.Fatal(http.ListenAndServe(":"+PORT, serveMux)) }()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	defer log.Fatal("Update ended?")
	for {
		select {
		case r := <-handler.request:
			update(db, r)
		case <-ticker.C:
			pullCitations(db, "2021-01-01") // TODO generate date
		}
	}
}
