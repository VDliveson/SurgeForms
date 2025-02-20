package database

import (
	"context"
	"log"

	"github.com/VDliveson/SurgeForms/go-producer/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client *mongo.Client
	Ctx    context.Context
}

var DBClient *mongo.Client

func ConnectDB(ctx context.Context) (*DB, error) {
	mongoURI := utils.GetEnv("MONGODB_URI", "")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	var err error

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MongoDB successfully")
	db := &DB{Client: client, Ctx: ctx}
	return db, nil
}

func (db *DB) GetCollection(collectionName string, databaseName string) *mongo.Collection {
	collection := db.Client.Database(databaseName).Collection(collectionName)
	return collection
}
