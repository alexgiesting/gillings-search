package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alexgiesting/gillings-search/paths"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	db          *mongo.Database
	client      *mongo.Client
	meta        *Collection
	Departments *Collection
	Faculty     *Collection
	Citations   *Collection
	Themes      *Collection
}

func (conn *Connection) Disconnect(ctx context.Context) error {
	return conn.client.Disconnect(ctx)
}

func (conn *Connection) Clear(ctx context.Context) error {
	return conn.meta.Drop(ctx)
}

func Connect() *Connection {
	DB_HOST := os.Getenv(paths.ENV_MONGODB_HOST)
	DB_PORT := os.Getenv(paths.ENV_MONGODB_PORT)
	DB_ADMIN_PASSWORD := os.Getenv(paths.ENV_MONGODB_ADMIN_PASSWORD)
	DB_NAME := os.Getenv(paths.ENV_MONGODB_NAME)

	DB_CREDENTIALS := ""
	if DB_ADMIN_PASSWORD != "" {
		DB_CREDENTIALS = fmt.Sprintf("admin:%s@", DB_ADMIN_PASSWORD)
	}
	DB_URI := fmt.Sprintf("mongodb://%s%s:%s/", DB_CREDENTIALS, DB_HOST, DB_PORT)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DB_URI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v\n", err)
	}
	db := client.Database(DB_NAME)

	return &Connection{
		db:          db,
		client:      client,
		meta:        &Collection{mongo: db.Collection(META)},
		Departments: &Collection{mongo: db.Collection(DEPARTMENTS)},
		Faculty:     &Collection{mongo: db.Collection(FACULTY)},
		Citations:   &Collection{mongo: db.Collection(CITATIONS)},
		Themes:      &Collection{mongo: db.Collection(THEMES)},
	}
}
