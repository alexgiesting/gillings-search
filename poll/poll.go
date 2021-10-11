package poll

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/alexgiesting/gillings-search/database"
	"github.com/alexgiesting/gillings-search/paths"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO not all fields are retrieved unless using a subscriber IP address
type ScopusResult struct {
	Results struct {
		Citations []struct {
			SID          string `json:"dc:identifier"`
			Title        string `json:"dc:title"`
			Author       string `json:"dc:creator"`
			PubType      string `json:"prism:aggregationType"`
			PubName      string `json:"prism:publicationName"`
			SubType      string `json:"subtypeDescription"`
			Volume       string `json:"prism:volume"`
			Pages        string `json:"prism:pageRange"`
			Date         string `json:"prism:coverDisplayDate"`
			ISODate      string `json:"prism:coverDate"`
			DOI          string `json:"prism:doi"`
			Abstract     string `json:"dc:description"`
			CitedByCount string `json:"citedby-count"`
			Keywords     string `json:"authkeywords"`
			Authors      []struct {
				SID            string `json:"authid"`
				Name           string `json:"authname"`
				GivenName      string `json:"given-name"`
				Surname        string `json:"surname"`
				Initials       string `json:"initials"`
				AffiliationSID string `json:"afid"`
			} `json:"author"`
			Affiliations []struct {
				SID     string `json:"afid"`
				Name    string `json:"affilname"`
				City    string `json:"affiliation-city"`
				Country string `json:"affiliation-country"`
			} `json:"affiliation"`
		} `json:"entry"`
		Count string `json:"opensearch:totalResults"`
	} `json:"search-results"`
}

func addCitations(db *mongo.Database, apiKey string) {
	cursor, err := db.Collection(database.FACULTY).Find(context.TODO(), bson.D{}, options.Find().SetProjection(bson.M{"sid": 1}))
	if err != nil {
		log.Fatal(err)
	}
	sidLists := make([]struct {
		SIDs []string `bson:"sid"`
	}, 0, cursor.RemainingBatchLength())
	err = cursor.All(context.TODO(), &sidLists)
	if err != nil {
		log.Fatal(err)
	}
	for _, sidList := range sidLists {
		for _, sid := range sidList.SIDs {
			query := url.QueryEscape(fmt.Sprintf("AU-ID(%s)", sid))
			url := fmt.Sprintf("https://api.elsevier.com/content/search/scopus?query=%s", query)
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatal(err)
			}
			request.Header.Set("Accept", "application/json")
			request.Header.Set("X-ELS-APIKey", apiKey)
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
			//addCitation(db, result) // TODO use chan
			log.Printf(result.Results.Citations[0].Title)
		}
		if true {
			return
		} // TODO
	}

	// query, err := json.Marshal(struct {
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
	// }{
	// 	Accept: "application/json",
	// 	View:   "COMPLETE",
	// })
	// if err != nil {
	// }

}

func Main() {
	client, db := database.Connect()
	defer client.Disconnect(context.TODO())
	database.Init(db)

	defer log.Fatal("Poll ended?")
	for {
		apiKey, present := os.LookupEnv(paths.ENV_SCOPUS_API_KEY)
		if !present {
			log.Fatal("Scopus API key missing")
		}
		addCitations(db, apiKey)
		time.Sleep(24 * time.Hour)
	}
}
