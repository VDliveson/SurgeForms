package utils

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client

func ConnectDB() error {
	mongoURI := GetEnv("MONGODB_URI", "")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	var err error

	log.Println("Connecting to MongoDB...")
	DBClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	//ping the database
	err = DBClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	log.Println("Connected to MongoDB successfully")
	return nil
}

func GetCollection(client *mongo.Client, collectionName string, databaseName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}
