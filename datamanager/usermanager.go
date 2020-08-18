package datamanager

import (
	"context"
	"github.com/google/uuid"
	"github.com/oprogramador/user-service-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userManagerStruct struct {
	usersCollection *mongo.Collection
	ctx             context.Context
}

func NewUserManager(usersCollection *mongo.Collection, ctx context.Context) *userManagerStruct {
	this := userManagerStruct{usersCollection, ctx}
	return &this
}

func (this *userManagerStruct) Save(user *models.User) error {
	if user.UserID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		user.UserID = id.String()
	}
	bsonData := bson.M{
		"active":  user.Active,
		"name":    user.Name,
		"user_id": user.UserID,
	}
	_, err := this.usersCollection.InsertOne(this.ctx, bsonData)
	return err
}

func (this *userManagerStruct) Delete(id string) error {
	_, err := this.usersCollection.DeleteOne(this.ctx, bson.D{
		{Key: "user_id", Value: id},
	})
	return err
}

func (this *userManagerStruct) Find(params ...map[string]interface{}) ([]models.User, error) {
	query := bson.M{}
	if len(params) > 0 {
		query = bson.M{"active": params[0]["active"].(bool)}
	}
	cursor, err := this.usersCollection.Find(this.ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(this.ctx)
	users := []models.User{}
	for cursor.Next(this.ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		user := models.User{UserID: result["user_id"].(string), Name: result["name"].(string), Active: result["active"].(bool)}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
