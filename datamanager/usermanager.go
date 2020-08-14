package datamanager

import (
	"context"
	"github.com/google/uuid"
	"github.com/oprogramador/user-service-golang/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userManager struct {
	usersCollection *mongo.Collection
	ctx             context.Context
}

func NewUserManager(usersCollection *mongo.Collection, ctx context.Context) *userManager {
	this := userManager{usersCollection, ctx}
	return &this
}

func (this *userManager) Save(user *models.User) error {
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

func (this *userManager) Delete(id string) error {
	_, err := this.usersCollection.DeleteOne(this.ctx, bson.D{
		{Key: "user_id", Value: id},
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
		user := models.User{UserID: result["user_id"].(string), Name: result["name"].(string), Active: result["active"].(bool)}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
