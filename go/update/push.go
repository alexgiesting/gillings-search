package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/alexgiesting/gillings-search/go/database"
	"github.com/alexgiesting/gillings-search/go/paths"
)

type Searchable struct {
	ID      string   `json:"id"`
	Text    string   `json:"_text_"`
	Title   string   `json:"title"`
	PubType string   `json:"pubtype"`
	PubName string   `json:"pubname"`
	Date    string   `json:"date"`
	CitedBy int      `json:"citedby"`
	Author  []string `json:"author"`
	SID     []string `json:"sid"`
}

func makeSearchable(citation database.Citation) Searchable {
	log.Print("a")
	authors := make([]string, len(citation.Authors))
	log.Print("b")
	sids := make([]string, len(citation.Authors))
	log.Print("c")
	for j, author := range citation.Authors {
		log.Print("d")
		authors[j] = fmt.Sprintf("%s %s %s", author.GivenName, author.Initials, author.Surname) // TODO
		log.Print("e")
		sids[j] = author.SID
		log.Print("f")
	}
	log.Print("g")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Printf("%d / %d", memStats.Alloc, memStats.HeapSys)
	return Searchable{
		ID:      citation.EID,
		Text:    fmt.Sprintf("%s %s %s", citation.Title, citation.Abstract, strings.Join(citation.Keywords, " ")), // TODO
		Title:   citation.Title,
		PubType: fmt.Sprintf("%s %s", citation.PubType, citation.SubType),
		PubName: fmt.Sprintf("%s %s %s", citation.PubName, citation.Volume, citation.Issue),
		Date:    fmt.Sprintf("%sT00:00:00Z", citation.ISODate),
		CitedBy: citation.CitedByCount,
		Author:  authors,
		SID:     sids,
	}
}

func pushCitations(db *database.Connection) {
	log.Print("pushing")

	var citations []database.Citation
	err := db.Citations.Decode(&citations)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("citations", len(citations))

	docs := make(chan []byte)
	i := 0
	j := 0
	go func() {
		for _, citation := range citations {
			doc, err := json.Marshal(makeSearchable(citation))
			if err != nil {
				log.Fatal(err)
			}
			docs <- doc
			i++
			j++
		}
		close(docs)
	}()

	host, _ := os.LookupEnv(paths.ENV_SOLR_HOST)
	port, _ := os.LookupEnv(paths.ENV_SOLR_PORT)
	url := fmt.Sprintf("http://%s:%s/solr/citations/update?", host, port)
	log.Print(url)
	body := make([]byte, 0, 64<<10)
	body = append(body, '[')
	for {
		doc, ok := <-docs
		if ok && len(body)+len(doc)+1 < cap(body) {
			body = append(body, doc...)
			body = append(body, ',')
			continue
		}

		body[len(body)-1] = ']'
		log.Print(i, len(body)) // TODO
		params := "overwrite=false"
		if j > 5000 {
			j = 0
			params = params + "&commit=true"
		}
		request, err := http.NewRequest("POST", url+params, bytes.NewBuffer(body))
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatal(i, err)
		}

		if !ok {
			break
		}
		body = append(body[:1], doc...)
		body = append(body, ',')
	}
	log.Print("Done pushing to Solr")
}
