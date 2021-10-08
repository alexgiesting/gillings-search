package adder

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func Main() {
	client, db := db_connect()
	defer client.Disconnect(context.TODO())
	db_init(db)

	go db_add_citations(db)
	select {}
}

func db_connect() (*mongo.Client, *mongo.Database) { // TODO extract to package??
	return nil, nil
}

func db_init(db *mongo.Database) {} // TODO should this really be in both places??

func db_add_citations(db *mongo.Database) {
	for {
		time.Sleep(time.Hour) // TODO
	}
}
