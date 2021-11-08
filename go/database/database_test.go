package database

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestConnect(t *testing.T) {
	db := Connect()
	db.Faculty.Drop(context.TODO())

	var documents []bson.D
	db.Faculty.Decode(&documents)
	assert(t, len(documents) == 0,
		"Expected collection to be empty after calling `Drop(ctx)`")
}

func assert(t *testing.T, condition bool, logValues ...interface{}) {
	if !condition {
		if len(logValues) > 0 {
			formatString := logValues[0].(string)
			if len(logValues) > 1 {
				logValues = logValues[1:]
			} else {
				logValues = nil
			}
			t.Logf(formatString, logValues...)
		}
		t.FailNow()
	}
}

func TestCollection(t *testing.T) {
	db := Connect()
	db.Faculty.Drop(context.TODO())
	db.Faculty.Insert(bson.M{
		"name":       "Jane",
		"age":        30,
		"user_color": "cyan",
	})
	db.Faculty.Insert(bson.M{
		"name":       "Philip",
		"age":        27,
		"user_color": "dark grey",
	})

	var ages []bson.M
	db.Faculty.Filter("name", "Jane").Project("age").Decode(&ages)
	assert(t, len(ages) == 1 && ages[0]["age"] == int32(30), // TODO fix this interface
		"Got %#v, but expected %#v", ages, []int32{30})

	var names []bson.M
	db.Faculty.Filter("user_color", bson.M{"$regex": `\bdark\b`}).Project("name").Decode(&names)
	assert(t, len(names) == 1 && names[0]["name"] == "Philip", // TODO fix this interface
		"Got %#v, but expected %#v", names, []string{"Philip"})

	type Result struct {
		Name  string `bson:"name"`
		Age   int    `bson:"age"`
		Color string `bson:"user_color"`
	}
	var results []Result
	db.Faculty.Filter("name", bson.M{"$in": []string{"Veronica", "Philip", "Simon", "Jane"}}).Decode(&results)
	assert(t, len(results) == 2, "Expected %#v to have two items", results)
	jane, philip := results[0], results[1]
	assert(t, jane.Name == "Jane" && jane.Age == 30 && jane.Color == "cyan",
		"Got %#v, but expected %#v", jane, Result{"Jane", 30, "cyan"})
	assert(t, philip.Name == "Philip" && philip.Age == 27 && philip.Color == "dark grey",
		"Got %#v, but expected %#v", philip, Result{"Philip", 27, "dark grey"})
}

func TestLoaders(t *testing.T) {
	db := Connect()

	db.Faculty.Drop(context.TODO())
	db.LoadFaculty(bytes.NewBuffer([]byte(
		`Name,Department,Title,Scopus ID,E-mail,Strengths
Ken B Donnel,mat,Clinical Associate Professor,3141592,ken@email.com,
Jennifer Tell,sci,Distinguished Professor,"98520830,58098932,41089432",jennifer@email.com,"ThemeA:Subtheme1,ThemeA:Subtheme2"`,
	)))
	var faculty []Faculty
	db.Faculty.Decode(&faculty)
	assert(t, len(faculty) == 2, "Expected to find two faculty documents after load")
	ken, jennifer := faculty[0], faculty[1]
	assert(t, ken.Name == "Ken B Donnel" &&
		ken.Department == "mat" &&
		ken.Title == "Clinical Associate Professor" &&
		len(ken.SID) == 1 && ken.SID[0] == "3141592" &&
		ken.Email == "ken@email.com",
		"Incorrect faculty data for Ken: %#v", ken,
	)
	assert(t, jennifer.Name == "Jennifer Tell" &&
		jennifer.Department == "sci" &&
		jennifer.Title == "Distinguished Professor" &&
		len(jennifer.SID) == 3 && jennifer.SID[0] == "98520830" && jennifer.SID[1] == "58098932" && jennifer.SID[2] == "41089432" &&
		jennifer.Email == "jennifer@email.com" &&
		len(jennifer.Strengths) == 2 && jennifer.Strengths[0].ThemeAbbr == "ThemeA" && jennifer.Strengths[0].SubThemeAbbr == "Subtheme1" &&
		jennifer.Strengths[1].ThemeAbbr == "ThemeA" && jennifer.Strengths[1].SubThemeAbbr == "Subtheme2",
		"Incorrect faculty data for Jennifer: %#v", jennifer,
	)

	db.Themes.Drop(context.TODO())
	db.LoadThemes(bytes.NewBuffer([]byte(
		`<?xml version="1.0" ?>
	<themes>
	
	<theme name="Research Theme A" abbr="A">
		<subtheme name="Subtheme 1" abbr="1">
			<definition>The first subtheme of Theme A!</definition>
			<keywords><kw>keyword (A.1.i)</kw><kw>keyword (A.1.ii)</kw><kw>keyword (A.1.iii)</kw></keywords>
		</subtheme>
		<subtheme name="Subtheme 2" abbr="2">
			<definition>The second subtheme of Theme A!</definition>
			<keywords><kw>keyword (A.2.i)</kw><kw>keyword (A.2.ii)</kw><kw>keyword (A.2.iii)</kw></keywords>
		</subtheme>
	</theme>
	
	<theme name="Research Theme B" abbr="B">
		<subtheme name="Subtheme 1" abbr="1"><definition>The first subtheme of Theme B!</definition><keywords></keywords></subtheme>
		<subtheme name="Subtheme 2" abbr="2"><definition>The second subtheme of Theme B!</definition><keywords></keywords></subtheme>
	</theme>

	<theme name="Research Theme C" abbr="C"></theme>
	
	</themes>
	`)))
	var themes []Theme
	db.Themes.Decode(&themes)
	assert(t, len(themes) == 3, "Expected to find three theme documents after load")
	themeA, themeB, themeC := themes[0], themes[1], themes[2]
	assert(t, themeA.Name == "Research Theme A" && themeA.Abbr == "A" &&
		len(themeA.SubThemes) == 2 &&
		themeA.SubThemes[0].Name == "Subtheme 1" && themeA.SubThemes[0].Abbr == "1" &&
		themeA.SubThemes[0].Description == "The first subtheme of Theme A!" &&
		strings.Join(themeA.SubThemes[0].Keywords, "|") == "keyword (A.1.i)|keyword (A.1.ii)|keyword (A.1.iii)" &&
		themeA.SubThemes[1].Name == "Subtheme 2" && themeA.SubThemes[1].Abbr == "2" &&
		themeA.SubThemes[1].Description == "The second subtheme of Theme A!" &&
		strings.Join(themeA.SubThemes[1].Keywords, "|") == "keyword (A.2.i)|keyword (A.2.ii)|keyword (A.2.iii)",
		"Incorrect theme data for Theme A: %+v", themeA,
	)
	assert(t, themeB.Name == "Research Theme B" && themeB.Abbr == "B" &&
		len(themeB.SubThemes) == 2 &&
		themeB.SubThemes[0].Name == "Subtheme 1" && themeB.SubThemes[0].Abbr == "1" &&
		themeB.SubThemes[0].Description == "The first subtheme of Theme B!" &&
		len(themeB.SubThemes[0].Keywords) == 0 &&
		themeB.SubThemes[1].Name == "Subtheme 2" && themeB.SubThemes[1].Abbr == "2" &&
		themeB.SubThemes[1].Description == "The second subtheme of Theme B!" &&
		len(themeB.SubThemes[1].Keywords) == 0,
		"Incorrect theme data for Theme B: %+v", themeB,
	)
	assert(t, themeC.Name == "Research Theme C" && themeC.Abbr == "C" &&
		len(themeC.SubThemes) == 0,
		"Incorrect theme data for Theme C: %+v", themeC,
	)
}
