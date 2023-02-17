package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func ConnectMongo(ctx context.Context, url string) *mongo.Client {
	clientOptions := options.Client()
	clientOptions.ApplyURI(url)
	clientOptions.SetAppName("beer")

	var err error
	if err = clientOptions.Validate(); err != nil {
		log.Fatalf("Mongo Client Options are invalid: %v\n", err)
	}

	mongoClient, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to Mongo Database! => %v\n", err)
	}

	if err = mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalln("Failed to ping Mongo Database!")
	}

	return mongoClient
}

func CloseMongoClient(ctx context.Context) {
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to close Mongo! => %v\n", err)
	}
}
