package datamanager

import (
	"context"
	"github.com/oprogramador/user-service-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userManager struct {
	usersCollection *mongo.Collection
	ctx             context.Context
}

func New(usersCollection *mongo.Collection, ctx context.Context) *userManager {
	this := userManager{usersCollection, ctx}
	return &this
}

func (this *userManager) Save(user *models.User) error {
	data, err := this.usersCollection.InsertOne(this.ctx, bson.D{
		{Key: "Active", Value: user.Active},
		{Key: "Name", Value: user.Name},
	})
	user.UserID = data.InsertedID.(primitive.ObjectID).Hex()
	return err
}

func (this *userManager) Delete(id string) error {
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = this.usersCollection.DeleteOne(this.ctx, bson.D{
		{Key: "_id", Value: idPrimitive},
	})
	return err
}

func (this *userManager) FindAll() ([]models.User, error) {
	cursor, err := this.usersCollection.Find(this.ctx, bson.D{})
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
		user := models.User{UserID: result["_id"].(primitive.ObjectID).Hex(), Name: result["Name"].(string), Active: result["Active"].(bool)}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
