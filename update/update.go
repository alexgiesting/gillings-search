package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alexgiesting/gillings-search/database"
	"github.com/alexgiesting/gillings-search/paths"
)

// type ScopusQuery struct {
// 	Accept           string `json:"httpAccept"` // application/json, application/atom+xml, application/xml
// 	AccessToken      string `json:"access_token"`
// 	InstitutionToken string `json:"insttoken"`
// 	APIKey           string `json:"apiKey"`
// 	RequestID        string `json:"reqId"`
// 	ResourceVersion  string `json:"ver"` // facetexpand, new
// 	Query            string `json:"query"`
// 	View             string `json:"view"` // STANDARD, COMPLETE
// 	SuppressNavLinks bool   `json:"suppressNavLinks"`
// 	Year             string `json:"date"`
// 	Offset           uint   `json:"start"`
// 	Count            uint   `json:"count"`
// 	Sort             string `json:"sort"`    // artnum, citedby-count, coverDate, creator, orig-load-date, pagecount, pagefirst, pageRange, publicationName, pubyear, relevancy, volume
// 	Content          string `json:"content"` // core, dummy, all
// 	Subject          string `json:"subj"`
// 	UseAuthorAlias   bool   `json:"alias"`
// 	Cursor           string `json:"cursor"`
// 	Facets           string `json:"facets"`
// }

type ScopusResult struct {
	Results struct {
		Citations []Entry `json:"entry"`
		Count     string  `json:"opensearch:totalResults"`
	} `json:"search-results"`
}

type Entry struct {
	EID          string        `json:"eid"`
	Title        string        `json:"dc:title"`
	Author       string        `json:"dc:creator"`
	PubType      string        `json:"prism:aggregationType"`
	PubName      string        `json:"prism:publicationName"`
	SubType      string        `json:"subtypeDescription"`
	Volume       string        `json:"prism:volume"`
	Pages        string        `json:"prism:pageRange"`
	Date         string        `json:"prism:coverDisplayDate"`
	ISODate      string        `json:"prism:coverDate"`
	DOI          string        `json:"prism:doi"`
	Abstract     string        `json:"dc:description"`
	CitedByCount string        `json:"citedby-count"`
	Keywords     string        `json:"authkeywords"`
	Authors      []EntryAuthor `json:"author"`
	Affiliations []struct {
		SID     string `json:"afid"`
		Name    string `json:"affilname"`
		City    string `json:"affiliation-city"`
		Country string `json:"affiliation-country"`
	} `json:"affiliation"`
}

type EntryAuthor struct {
	SID          string `json:"authid"`
	Name         string `json:"authname"`
	GivenName    string `json:"given-name"`
	Surname      string `json:"surname"`
	Initials     string `json:"initials"`
	Affiliations []struct {
		SID string `json:"$"`
	} `json:"afid"`
}

func addCitations(db *database.Connection) {
	apiKey, present := os.LookupEnv(paths.ENV_SCOPUS_API_KEY)
	if !present {
		log.Fatal("Scopus API key missing")
	}
	apiClient, _ := os.LookupEnv(paths.ENV_SCOPUS_CLIENT_ADDRESS)

	for _, sid := range getSIDs(db) {
		result := queryScopus(sid, apiKey, apiClient)
		for _, entry := range result.Results.Citations {
			exists, err := db.Citations.Filter("eid", entry.EID).Check()
			if err != nil {
				log.Fatal(err)
			}
			if !exists {
				addCitation(db, &entry) // TODO use chan instead?
			}
		}
		break
		// TODO make testing version
		// TODO make a version that only adds recent results
		// TODO make sure process can recover from interruptions
	}
}

func getSIDs(db *database.Connection) []string {
	var sidLists []struct{ SID []string }
	err := db.Faculty.Project("sid").Decode(&sidLists)
	if err != nil {
		log.Fatal(err)
	}

	var sids []string
	for _, sidList := range sidLists {
		sids = append(sids, sidList.SID...)
	}
	return sids
}

func queryScopus(sid string, apiKey string, apiClient string) ScopusResult {
	query := url.QueryEscape(fmt.Sprintf("AU-ID(%s)", sid))
	url := fmt.Sprintf("https://api.elsevier.com/content/search/scopus?query=%s&view=COMPLETE", query)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-ELS-APIKey", apiKey)
	if apiClient != "" {
		request.Header.Set("X-Forwarded-For", apiClient)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result ScopusResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func addCitation(db *database.Connection, entry *Entry) {
	err := db.Citations.Insert(database.Citation{
		Title:        entry.Title,
		PubType:      entry.PubType,
		PubName:      entry.PubName,
		SubType:      entry.SubType,
		Volume:       entry.Volume,
		Pages:        entry.Pages,
		Date:         entry.Date,
		ISODate:      entry.ISODate,
		DOI:          entry.DOI,
		Abstract:     entry.Abstract,
		CitedByCount: entry.getCitedByCount(),
		Keywords:     entry.getKeywords(),
		EID:          entry.EID,
		Authors:      entry.getAuthors(),
		Affiliations: entry.getAffiliations(),
		Status:       database.STATUS_UNCONFIRMED,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (entry *Entry) getCitedByCount() int {
	citedByCount, err := strconv.Atoi(entry.CitedByCount)
	if err != nil {
		log.Fatal(err)
	}
	return citedByCount
}

func (entry *Entry) getKeywords() []string {
	if entry.Keywords == "" {
		return []string{}
	} else {
		return strings.Split(entry.Keywords, " | ")
	}
}

func (entry *Entry) getAuthors() []database.Author {
	authors := make([]database.Author, len(entry.Authors))
	for i, author := range entry.Authors {
		affiliations := make([]string, len(author.Affiliations))
		for i, affiliation := range author.Affiliations {
			affiliations[i] = affiliation.SID
		}
		authors[i] = database.Author{
			GivenName: author.GivenName,
			Surname:   author.Surname,
			SID:       author.SID,
			AffilIDs:  affiliations,
		}
	}
	return authors
}

func (entry *Entry) getAffiliations() []database.Affiliation {
	affiliations := make([]database.Affiliation, len(entry.Affiliations))
	for i, affiliation := range entry.Affiliations {
		affiliations[i] = database.Affiliation{
			SID:     affiliation.SID,
			Name:    affiliation.Name,
			City:    affiliation.City,
			Country: affiliation.Country,
		}
	}
	return affiliations
}

type QueryHandler struct {
	request chan Request
}

type Request uint

const (
	PULL Request = iota
	INITIALIZE
	RESET
)

func (handler *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query["key"][0] != os.Getenv(paths.ENV_UPDATE_KEY) {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	path := r.URL.Path[len(paths.PATH_UPDATE):]
	path = strings.TrimRight(path, "/")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, "command received")

	switch path {
	case "pull":
		handler.request <- PULL
	case "init":
		handler.request <- INITIALIZE
	case "reset":
		handler.request <- RESET
	default:
		log.Printf("Invalid request `%s` received by `update`\n", path)
	}
}

func Main() {
	db := database.Connect()
	defer db.Disconnect(context.TODO())
	db.Init()

	serveMux := http.NewServeMux()
	handler := QueryHandler{make(chan Request)}
	serveMux.Handle(paths.PATH_UPDATE, &handler)

	PORT := os.Getenv(paths.ENV_UPDATE_PORT)
	log.Printf("Running server on %s\n", PORT)
	go func() { log.Fatal(http.ListenAndServe(PORT, serveMux)) }()

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	defer log.Fatal("Update ended?")
	for {
		select {
		case r := <-handler.request:
			switch r {
			case PULL:
				addCitations(db)
			case INITIALIZE:
				db.Init()
			case RESET:
				db.Clear(context.TODO())
				db.Init()
			}
		case <-ticker.C:
			addCitations(db)
		}
	}
}
