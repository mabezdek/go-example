package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_DSN"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	return client
}

var Database *mongo.Client = Connect()

func GetCollection(collectionName string) *mongo.Collection {
	return Database.Database("rubumo").Collection(collectionName)
}

func GetApplicationCollection(applicationId string, collectionName string) *mongo.Collection {
	return Database.Database(applicationId).Collection(collectionName)
}
