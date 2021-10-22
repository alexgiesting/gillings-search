package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO refactor...

type Collection struct {
	mongo *mongo.Collection
}

func (collection *Collection) Drop(ctx context.Context) error {
	return collection.mongo.Drop(ctx)
}

func (collection *Collection) Insert(item interface{}) error {
	document, err := bson.Marshal(item)
	if err != nil {
		return err
	}
	_, err = collection.mongo.InsertOne(context.TODO(), document)
	return err
}

func (collection *Collection) Decode(results interface{}) error {
	cursor, err := collection.mongo.Find(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	return cursor.All(context.TODO(), results)
}

type Search struct {
	collection *mongo.Collection
	filter     []bson.D
	or         [][]bson.D
	project    bson.D
}
type Projection Search

func (collection *Collection) Search() *Search {
	return newSearch(collection)
}

func (collection *Collection) Filter(key string, filter interface{}) *Search {
	return newSearch(collection).Filter(key, filter)
}

func (collection *Collection) Project(includeFields ...string) *Projection {
	return newSearch(collection).Project(includeFields...)
}

func newSearch(collection *Collection) *Search {
	filter := []bson.D{}
	or := [][]bson.D{}
	project := bson.D{}
	return &Search{collection.mongo, filter, or, project}
}

func (search *Search) makeFilter() bson.D {
	filter := search.filter
	if len(search.or) > 0 {
		orFilter := make([]bson.D, len(search.or))
		for i, alternative := range search.or {
			orFilter[i] = bson.D{{Key: "$and", Value: alternative}}
		}
		filter = append(filter, bson.D{{Key: "$or", Value: orFilter}})
	}
	return bson.D{{Key: "$and", Value: filter}}
}

func (search *Search) Filter(key string, filter interface{}) *Search {
	search.filter = append(search.filter, bson.D{{Key: key, Value: filter}})
	return search
}

func (search *Search) Check() (bool, error) {
	result := search.collection.FindOne(context.TODO(), search.filter)
	if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if result.Err() != nil {
		return false, result.Err()
	} else {
		return true, nil
	}
}

func (search *Search) Decode(results interface{}) error {
	cursor, err := search.collection.Find(context.TODO(), search.makeFilter())
	if err != nil {
		return err
	}
	return cursor.All(context.TODO(), results)
}

func (search *Search) Project(fields ...string) *Projection {
	for _, field := range fields {
		search.project = append(search.project, bson.E{Key: field, Value: 1})
	}
	return (*Projection)(search)
}

func (projection *Projection) Decode(results interface{}) error {
	filter := (*Search)(projection).makeFilter()
	option := options.Find().SetProjection(projection.project)
	cursor, err := projection.collection.Find(context.TODO(), filter, option)
	if err != nil {
		return err
	}
	return cursor.All(context.TODO(), results)
}

type Alternative struct {
	search *Search
	index  int
}

func (search *Search) Alternative() *Alternative {
	index := len(search.or)
	search.or = append(search.or, []bson.D{})
	return &Alternative{search, index}
}

func (alt *Alternative) Filter(key string, filter interface{}) *Alternative {
	alt.search.or[alt.index] = append(alt.search.or[alt.index], bson.D{{Key: key, Value: filter}})
	return alt
}
