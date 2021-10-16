package database

import (
	"context"
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	META        = "__dbinfo__"
	DEPARTMENTS = "Departments"
	FACULTY     = "Faculty"
	CITATIONS   = "Citations"
	THEMES      = "Themes"
)

func (conn *Connection) Init() {
	db := conn.db
	colls := make(map[string]bool)
	names, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	for _, name := range names {
		colls[name] = true
	}

	if colls[META] {
		var dbInfo DatabaseInfo
		db.Collection(META).FindOne(context.TODO(), bson.D{}).Decode(&dbInfo)
		if dbInfo.Initialized {
			return
		}
		for _, coll := range dbInfo.UninitializedCollections {
			switch coll.Name {
			case CITATIONS:
				initCitations(db, coll.Recovery)
			case THEMES:
				initThemes(db)
			case FACULTY:
				initFaculty(db)
			}
		}
	} else {
		db.CreateCollection(context.TODO(), META)
		db.Collection(META).InsertOne(context.TODO(), DatabaseInfo{
			Initialized: false,
			UninitializedCollections: []struct {
				Name     string
				Recovery string
			}{
				{Name: FACULTY},
				{Name: THEMES},
				{Name: CITATIONS, Recovery: ""},
			},
		})
		initFaculty(db)
		initThemes(db)
		initCitations(db, "")
		_, err = db.Collection(META).UpdateOne(context.TODO(), bson.D{}, bson.M{
			"$set": bson.M{"initialized": "true"},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func drop(db *mongo.Database, collection string) {
	err := db.Collection(collection).Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	err = db.CreateCollection(context.TODO(), collection)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func pop(db *mongo.Database) {
	_, err := db.Collection(META).UpdateOne(context.TODO(), bson.D{}, bson.M{
		"$pop": bson.M{"uninitializedcollections": 1},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func initFaculty(db *mongo.Database) {
	drop(db, FACULTY)
	drop(db, DEPARTMENTS)

	facultyCSV, err := os.Open("data/faculty.csv")
	if err != nil {
		log.Fatal(err)
	}
	lines, err := csv.NewReader(facultyCSV).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	faculty := make([]interface{}, len(lines)-1)
	departments := make(map[string][]string)
	for r, row := range lines[1:] {
		fields := make(map[string]string)
		for c, label := range lines[0] {
			fields[label] = row[c]
		}
		name := strings.Split(fields["Name"], ", ")
		sid := strings.Split(fields["Scopus ID"], ",")
		faculty[r] = Faculty{
			GivenName: name[1],
			Surname:   name[0],
			Title:     fields["Title"],
			SID:       sid,
		}
		departments[fields["Department"]] = append(departments[fields["Department"]], sid...)
	}

	_, err = db.Collection(FACULTY).InsertMany(context.TODO(), faculty)
	if err != nil {
		log.Fatal(err)
	}
	for name, sids := range departments {
		_, err = db.Collection(FACULTY).InsertOne(context.TODO(), Department{
			Name: name, SIDs: sids,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	pop(db)
}

func initThemes(db *mongo.Database) {
	drop(db, THEMES)

	themesXML, err := os.Open("data/themes.xml")
	if err != nil {
		log.Fatal(err)
	}
	themesBytes, err := ioutil.ReadAll(themesXML)
	if err != nil {
		log.Fatal(err)
	}

	themes := make([]Theme, 0, 10)
	err = xml.Unmarshal(themesBytes, &themes)
	if err != nil {
		log.Fatal(err)
	}
	for _, theme := range themes {
		_, err = db.Collection(THEMES).InsertOne(context.TODO(), theme)
		if err != nil {
			log.Fatal(err)
		}
	}

	pop(db)
}

func initCitations(db *mongo.Database, recovery string) {
	drop(db, CITATIONS)
	// TODO use update?
	pop(db)
}
