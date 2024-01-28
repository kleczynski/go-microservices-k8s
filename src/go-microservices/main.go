package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/kleczynski/go-microservices-k8s/database"
	"github.com/kleczynski/go-microservices-k8s/router"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = database.CreateConnection()
	if err != nil {
		log.Fatalf("There was an error creating connection with MongoDB: %v\n", err)
	}
	defer database.Client.Disconnect(context.Background())
	database.PingDB(database.Client)
	if err := router.StartServer(":5555"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
