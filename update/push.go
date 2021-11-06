package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/alexgiesting/gillings-search/database"
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

func pushCitations(db *database.Connection) {
	var citations []database.Citation
	db.Citations.Decode(&citations)

	searchables := make([]Searchable, len(citations))
	for i, citation := range citations {
		authors := make([]string, len(citation.Authors))
		sids := make([]string, len(citation.Authors))
		for j, author := range citation.Authors {
			authors[j] = fmt.Sprintf("%s %s %s", author.GivenName, author.Initials, author.Surname) // TODO
			sids[j] = author.SID
		}
		searchables[i] = Searchable{
			ID:      citation.EID,
			Text:    fmt.Sprintf("%s %s %s", citation.Title, citation.Abstract, strings.Join(citation.Keywords, " ")), // TODO
			Title:   citation.Title,
			PubType: fmt.Sprintf("%s %s", citation.PubType, citation.SubType),
			PubName: fmt.Sprintf("%s %s %s", citation.PubName, citation.Volume, citation.Issue),
			Date:    fmt.Sprintf("%sT23:59:59Z", citation.ISODate),
			CitedBy: citation.CitedByCount,
			Author:  authors,
			SID:     sids,
		}
	}

	docs := make(chan []byte)
	i := 0
	j := 0
	go func() {
		for _, searchable := range searchables {
			doc, err := json.Marshal(searchable)
			if err != nil {
				log.Fatal(err)
			}
			docs <- doc
			i++
			j++
		}
		close(docs)
	}()

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
		log.Print(i, len(body))
		params := "overwrite=false"
		if j > 5000 {
			j = 0
			params = params + "&commit=true"
		}
		request, err := http.NewRequest("POST", "http://localhost:8983/solr/citations/update?"+params, bytes.NewBuffer(body))
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
}
