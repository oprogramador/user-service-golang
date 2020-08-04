package main

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func listUsers(ctx context.Context, usersCollection *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cursor, err := usersCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln(err)
		}
		defer cursor.Close(ctx)
		users := []interface{}{}
		for cursor.Next(ctx) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("users", result)
			users = append(users, result)
		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(users)
	}
}

func createUser(ctx context.Context, usersCollection *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cursor, err := usersCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln(err)
		}
		defer cursor.Close(ctx)
		reqBody, _ := ioutil.ReadAll(r.Body)
		var user User
		json.Unmarshal(reqBody, &user)
		log.Println("user", user)

		data, err := usersCollection.InsertOne(ctx, bson.D{
			{Key: "Active", Value: user.Active},
			{Key: "Name", Value: user.Name},
		})
		log.Println(data, err)
	}
}

func deleteUser(ctx context.Context, usersCollection *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query()["id"][0]
		log.Println("delete", id)

		idPrimitive, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatalln(err)
		}

		data, err := usersCollection.DeleteOne(ctx, bson.D{
			{Key: "_id", Value: idPrimitive},
		})
		log.Println(data, err)
	}
}

func handleRequests(port string, ctx context.Context, usersCollection *mongo.Collection) {
	http.HandleFunc("/users", listUsers(ctx, usersCollection))
	http.HandleFunc("/user", createUser(ctx, usersCollection))
	http.HandleFunc("/delete-user", deleteUser(ctx, usersCollection))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	port := "10000"
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	usersDatabase := client.Database("users")
	usersCollection := usersDatabase.Collection("users")

	handleRequests(port, ctx, usersCollection)
}
