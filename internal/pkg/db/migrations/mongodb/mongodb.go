package database

import (
	"context"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDbConfig *mongodb.Config

var MongoDbClient *mongo.Client

func InitMongoDB() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}

	MongoDbClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return MongoDbClient
}

func MigrateMongoDB() {

}

func CloseMongoDB() {
	if err := MongoDbClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
