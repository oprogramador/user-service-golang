package main

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"
)

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func listUsers(ctx context.Context, usersCollection *mongo.Collection) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		cursor, err := usersCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln(err)
		}
		defer cursor.Close(ctx)
		users := []User{}
		for cursor.Next(ctx) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatalln(err)
			}
			user := User{UserID: result["_id"].(primitive.ObjectID).Hex(), Name: result["Name"].(string), Active: result["Active"].(bool)}
			users = append(users, user)
		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
		ginContext.JSON(200, users)
	}
}

func createUser(ctx context.Context, usersCollection *mongo.Collection) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		cursor, err := usersCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln(err)
		}
		defer cursor.Close(ctx)
		reqBody, _ := ioutil.ReadAll(ginContext.Request.Body)
		var user User
		err = json.Unmarshal(reqBody, &user)
		if err != nil {
			re := regexp.MustCompile(`[A-Za-z.]* of type [A-Za-z]*`)
			ginContext.String(400, strings.ReplaceAll(string(re.Find([]byte(err.Error()))), "of type", "should be of type"))
			return
		}

		data, err := usersCollection.InsertOne(ctx, bson.D{
			{Key: "Active", Value: user.Active},
			{Key: "Name", Value: user.Name},
		})
		log.Println(data, err)
		user.UserID = data.InsertedID.(primitive.ObjectID).Hex()
		ginContext.JSON(201, user)
	}
}

func deleteUser(ctx context.Context, usersCollection *mongo.Collection) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		id := ginContext.Param("id")
		log.Println("delete", id)

		idPrimitive, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			ginContext.String(400, "invalid id")
			return
		}

		data, err := usersCollection.DeleteOne(ctx, bson.D{
			{Key: "_id", Value: idPrimitive},
		})
		log.Println(data, err)
		ginContext.String(204, "")
	}
}

func handleRequests(ctx context.Context, usersCollection *mongo.Collection) *gin.Engine {
	router := gin.Default()
	router.GET("/users", listUsers(ctx, usersCollection))
	router.POST("/user", createUser(ctx, usersCollection))
	router.DELETE("/user/:id", deleteUser(ctx, usersCollection))

	return router
}

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

	router := handleRequests(ctx, usersCollection)

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
