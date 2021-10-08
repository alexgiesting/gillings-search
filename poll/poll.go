package poll

import (
	"context"
	"time"

	"github.com/alexgiesting/gillings-search/database"
	"go.mongodb.org/mongo-driver/mongo"
)

func addCitations(db *mongo.Database) {
	for {
		time.Sleep(time.Hour) // TODO
	}
}

func Main() {
	client, db := database.Connect()
	defer client.Disconnect(context.TODO())
	database.Init(db)

	go addCitations(db)
	select {}
}
