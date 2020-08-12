package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func ListUsers(ctx context.Context, usersCollection *mongo.Collection) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		cursor, err := usersCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln(err)
		}
		defer cursor.Close(ctx)
		users := []models.User{}
		for cursor.Next(ctx) {
			var result bson.M
			err := cursor.Decode(&result)
			if err != nil {
				log.Fatalln(err)
			}
			user := models.User{UserID: result["_id"].(primitive.ObjectID).Hex(), Name: result["Name"].(string), Active: result["Active"].(bool)}
			users = append(users, user)
		}
		if err := cursor.Err(); err != nil {
			log.Fatal(err)
		}
		ginContext.JSON(200, users)
	}
}
