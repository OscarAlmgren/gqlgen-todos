package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoDbClient *mongo.Client

func InitMongoDB() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	MongoDbClient = client
}

func PingMongoDB() {
	if err := MongoDbClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println("ping mongo break")
		panic(err)
	}
	fmt.Println("MongoDB Ping OK")
}

func MigrateMongoDB() {

}

func CloseMongoDB() {
	if err := MongoDbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
