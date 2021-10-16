package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func (collection *Collection) Search() *Search {
	return &Search{
		collection: collection.mongo,
		filters:    make(map[string]interface{}),
	}
}

func (collection *Collection) Filter(key string, filter interface{}) *Search {
	return collection.Search().Filter(key, filter)
}

func (collection *Collection) Project(includeFields ...string) *Projection {
	return (&Search{collection: collection.mongo, filters: nil}).Project(includeFields...)
}

type Search struct {
	collection *mongo.Collection
	filters    map[string]interface{}
}

func (search *Search) Filter(key string, filter interface{}) *Search {
	search.filters[key] = filter
	return search
}

func (search *Search) Project(includeFields ...string) *Projection {
	mapFields := make(map[string]int)
	for _, field := range includeFields {
		mapFields[field] = 1
	}
	fields, _ := bson.Marshal(mapFields)
	return &Projection{search, fields}
}

func (search *Search) Decode(results interface{}) error {
	bsonFilter, err := bson.Marshal(search.filters)
	if err != nil {
		return err
	}
	cursor, err := search.collection.Find(context.TODO(), bsonFilter)
	if err != nil {
		return err
	}
	return cursor.All(context.TODO(), results)
}

func (search *Search) Check() (bool, error) {
	bsonFilter, err := bson.Marshal(search.filters)
	if err != nil {
		return false, err
	}
	result := search.collection.FindOne(context.TODO(), bsonFilter)
	if result.Err() == mongo.ErrNoDocuments {
		return false, nil
	} else if result.Err() != nil {
		return false, result.Err()
	} else {
		return true, nil
	}
}

type Projection struct {
	search *Search
	fields interface{}
}

func (projection *Projection) Decode(results interface{}) error {
	bsonFilter, err := bson.Marshal(projection.search.filters)
	if err != nil {
		return err
	}
	cursor, err := projection.search.collection.Find(
		context.TODO(),
		bsonFilter,
		options.Find().SetProjection(projection.fields),
	)
	if err != nil {
		return err
	}
	return cursor.All(context.TODO(), results)
}
