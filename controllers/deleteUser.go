package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func DeleteUser(ctx context.Context, usersCollection *mongo.Collection) func(ginContext *gin.Context) {
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
