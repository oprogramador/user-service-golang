package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/datamanager"
	"github.com/oprogramador/user-service-golang/routing"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func disconnect(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func setupServer() (*gin.Engine, context.CancelFunc, *mongo.Client, context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	usersDatabase := client.Database("users")
	usersCollection := usersDatabase.Collection("users")

	userManager := datamanager.New(usersCollection, ctx)
	router := routing.HandleRequests(userManager)

	return router, cancel, client, ctx
}

func main() {
	port := "10000"
	router, cancel, client, ctx := setupServer()
	defer cancel()
	defer disconnect(client, ctx)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatalln(err)
	}
}
