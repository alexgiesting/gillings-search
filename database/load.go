package database

import (
	"context"
	"encoding/csv"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

const (
	META        = "__dbinfo__"
	DEPARTMENTS = "Departments"
	FACULTY     = "Faculty"
	CITATIONS   = "Citations"
	THEMES      = "Themes"
)

func (db *Connection) LoadFaculty(facultyCSV io.Reader) {
	lines, err := csv.NewReader(facultyCSV).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	faculty := make([]interface{}, len(lines)-1)
	for r, row := range lines[1:] {
		fields := make(map[string]string)
		for c, label := range lines[0] {
			fields[label] = row[c]
		}

		name := strings.Split(fields["Name"], ", ")
		sid := strings.Split(fields["Scopus ID"], ",")
		strengths := getStrengths(fields["Strengths"])
		faculty[r] = Faculty{
			GivenName: name[1],
			Surname:   name[0],
			Title:     fields["Title"],
			SID:       sid,
			Email:     fields["E-mail"],
			Strengths: strengths,
		}
	}

	_, err = db.Faculty.mongo.InsertMany(context.TODO(), faculty)
	if err != nil {
		log.Fatal(err)
	}
}

func getStrengths(strengthsString string) []Strength {
	if strengthsString == "" {
		return []Strength{}
	}
	strengthsStrings := strings.Split(strengthsString, ",")
	strengths := make([]Strength, len(strengthsStrings))
	for i, strength := range strengthsStrings {
		parts := strings.Split(strength, ":")
		strengths[i].Theme = parts[0]
		strengths[i].SubTheme = parts[1]
	}
	return strengths
}

func (db *Connection) LoadThemes(themesXML io.Reader) {
	themesBytes, err := ioutil.ReadAll(themesXML)
	if err != nil {
		log.Fatal(err)
	}

	var themes struct {
		Themes []Theme `xml:"theme"`
	}
	err = xml.Unmarshal(themesBytes, &themes)
	if err != nil {
		log.Fatal(err)
	}
	for _, theme := range themes.Themes {
		_, err = db.Themes.mongo.InsertOne(context.TODO(), theme)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (db *Connection) LoadCitations(citationsJSON io.Reader) {
	// TODO
}
