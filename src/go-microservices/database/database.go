package database

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func CreateConnection() error {
	credentails := options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}
	connString := os.Getenv("MONGO_URI")
	clientOpts := options.Client().ApplyURI(connString).SetAuth(credentails)
	var err error
	Client, err = mongo.Connect(context.TODO(), clientOpts)
	return err
}

func PingDB(client *mongo.Client) error {
	db := client.Database(os.Getenv("MONGO_DATABASE"))
	result := db.RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}})
	if result.Err() != nil {
		return fmt.Errorf("ping command failed: %v", result.Err())
	}
	fmt.Printf("Pinged %s database. Successfully connected\n", os.Getenv("MONGO_DATABASE"))
	return nil
}
