package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/alexgiesting/gillings-search/database"
	"github.com/alexgiesting/gillings-search/paths"
)

func pullCitations(db *database.Connection, startDate string) {
	apiKey, present := os.LookupEnv(paths.ENV_SCOPUS_API_KEY)
	if !present {
		log.Fatal("Scopus API key missing")
	}
	apiClient, _ := os.LookupEnv(paths.ENV_SCOPUS_CLIENT_ADDRESS)

	limiter := make(chan int, 8)
	for _, sids := range getSIDs(db) {
		if len(sids) == 0 {
			continue
		}

		limiter <- 1
		go func(sids []string) {
			entries := queryScopus(sids, startDate, apiKey, apiClient)
			for _, entry := range entries {
				exists, err := db.Citations.Filter("eid", entry.EID).Check()
				if err != nil {
					log.Fatal(err)
				}
				if !exists {
					addCitation(db, &entry)
				}
			}
			<-limiter
		}(sids)
		// TODO make a version that only adds recent results
		// TODO make a version that alters records based on faculty changes
		// TODO make sure process can recover from interruptions
	}
}

func getSIDs(db *database.Connection) [][]string {
	var sidLists []struct{ SID []string }
	err := db.Faculty.Project("sid").Decode(&sidLists)
	if err != nil {
		log.Fatal(err)
	}

	sids := make([][]string, len(sidLists))
	for i, sidList := range sidLists {
		sids[i] = sidList.SID
	}
	return sids
}

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

func queryScopus(sids []string, startDate string, apiKey string, apiClient string) []Entry {
	// TODO use date to limit results
	// TODO use EID to limit results
	// TODO monitor rate limits, request ids, errors...
	// TODO add progress logging
	fields := make([]string, len(sids))
	for i, sid := range sids {
		fields[i] = fmt.Sprintf("AU-ID(%s)", sid)
	}
	query := url.QueryEscape(strings.Join(fields, " OR "))
	url := fmt.Sprintf("https://api.elsevier.com/content/search/scopus?query=%s&view=COMPLETE", query)

	var entries []Entry
	var start int = 0
	var count int
	for {
		startField := fmt.Sprintf("&start=%d", start)
		request, err := http.NewRequest("GET", url+startField, nil)
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

		if start == 0 {
			i, err := strconv.Atoi(result.Results.Count)
			if err != nil {
				log.Fatal(err, url+startField, response, string(body), result)
			}
			count = i
			entries = make([]Entry, 0, count)
		}
		entries = append(entries, result.Results.Citations...)
		start += len(result.Results.Citations)
		if start >= count {
			break
		}
	}

	return entries
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
