package poll

import (
	"context"
	"time"

	"github.com/alexgiesting/gillings-search/db"
)

func addCitations(db *db.Database) {
	for {
		time.Sleep(time.Hour) // TODO
	}
}

func Main() {
	client, db := db.Connect()
	defer client.Disconnect(context.TODO())
	db.Init()

	go addCitations(db)
	select {}
}
