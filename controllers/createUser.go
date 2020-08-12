package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/oprogramador/user-service-golang/models"
	"io/ioutil"
	"encoding/json"
	"log"
	"regexp"
	"strings"
)

func CreateUser(ctx context.Context, usersCollection *mongo.Collection) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		cursor, err := usersCollection.Find(ctx, bson.D{})
		if err != nil {
			log.Fatalln(err)
		}
		defer cursor.Close(ctx)
		reqBody, _ := ioutil.ReadAll(ginContext.Request.Body)
		var user models.User
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

