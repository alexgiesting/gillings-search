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

	fields := make(map[string]int)
	for c, label := range lines[0] {
		fields[label] = c
	}
	NAME := fields["Name"]
	TITLE := fields["Title"]
	SID := fields["Scopus ID"]
	EMAIL := fields["E-mail"]
	STRENGTHS := fields["Strengths"]

	faculty := make([]interface{}, len(lines)-1)
	for r, row := range lines[1:] {
		faculty[r] = Faculty{
			Name:      row[NAME],
			Title:     row[TITLE],
			SID:       getSIDs(row[SID]),
			Email:     row[EMAIL],
			Strengths: getStrengths(row[STRENGTHS]),
		}
	}

	_, err = db.Faculty.mongo.InsertMany(context.TODO(), faculty)
	if err != nil {
		log.Fatal(err)
	}
}

func getSIDs(sidsString string) []string {
	if sidsString == "" {
		return []string{}
	}
	return strings.Split(sidsString, ",")
}

func getStrengths(strengthsString string) []Strength {
	if strengthsString == "" {
		return []Strength{}
	}
	strengthsStrings := strings.Split(strengthsString, ",")
	strengths := make([]Strength, len(strengthsStrings))
	for i, strength := range strengthsStrings {
		parts := strings.Split(strength, ":")
		strengths[i].ThemeAbbr = parts[0]
		strengths[i].SubThemeAbbr = parts[1]
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
