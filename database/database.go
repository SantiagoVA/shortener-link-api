package database

import (
	"context"
	"shortener-app/functions"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(collection string) *mongo.Collection {
	URI := functions.GetEnv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))

	if err != nil {
		panic(err)
	}

	database := client.Database("shortener")

	return database.Collection(collection)
}
