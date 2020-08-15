package main

import (
	"context"
	"github.com/oprogramador/user-service-golang/datamanager"
	"github.com/oprogramador/user-service-golang/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

type UserManager interface {
	Save(user *models.User) error
	Delete(id string) error
	FindAll() ([]models.User, error)
}

var userManager UserManager
var ctx context.Context
var cancel context.CancelFunc
var usersCollection *mongo.Collection

func beforeEach() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	usersDatabase := client.Database("users")
	err = usersDatabase.Drop(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	usersCollection = usersDatabase.Collection("users")

	userManager = datamanager.NewUserManager(usersCollection, ctx)
}

func afterEach() {
	defer cancel()
}

func TestSavingWithCustomId(t *testing.T) {
	beforeEach()
	defer afterEach()
	id := "c81dc894-3d59-4f02-b22b-d4ad2cba0610"
	user := models.User{Name: "Alan", Active: true, UserID: id}
	err := userManager.Save(&user)

	assert.Nil(t, err)
	assert.Equal(t, "Alan", user.Name)
	assert.Equal(t, true, user.Active)
	assert.Equal(t, id, user.UserID)

	var saved models.User
	err = usersCollection.FindOne(ctx, bson.M{"user_id": id}).Decode(&saved)
	assert.Nil(t, err)
	assert.Equal(t, user, saved)
}

func TestSavingWithoutCustomId(t *testing.T) {
	beforeEach()
	defer afterEach()
	user := models.User{Name: "Alan", Active: true}
	err := userManager.Save(&user)

	assert.Nil(t, err)
	assert.Equal(t, "Alan", user.Name)
	assert.Equal(t, true, user.Active)
	assert.Equal(t, 36, len(user.UserID))

	var saved models.User
	err = usersCollection.FindOne(ctx, bson.M{"user_id": user.UserID}).Decode(&saved)
	assert.Nil(t, err)
	assert.Equal(t, user, saved)
}
