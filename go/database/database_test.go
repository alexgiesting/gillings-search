package database

import (
	"context"
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
	// TODO
}
