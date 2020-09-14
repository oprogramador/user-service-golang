package datamanager

import (
	"context"
	. "github.com/franela/goblin"
	"github.com/oprogramador/user-service-golang/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
	"time"
)

var userManager *userManagerStruct
var ctx context.Context
var cancel context.CancelFunc
var usersCollection *mongo.Collection

func Test(t *testing.T) {
	g := Goblin(t)
	g.Describe("UserManager", func() {
		g.BeforeEach(func() {
			client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
			if err != nil {
				log.Fatalln(err)
			}
			ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
			err = client.Connect(ctx)
			if err != nil {
				log.Fatalln(err)
			}

			usersDatabase := client.Database("users")
			err = usersDatabase.Drop(ctx)
			if err != nil {
				log.Fatalln(err)
			}
			usersCollection = usersDatabase.Collection("users")

			userManager = NewUserManager(usersCollection, ctx)
		})

		g.AfterEach(func() {
			defer cancel()
		})

		g.It("saves with custom id", func() {
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
		})

		g.It("saves without custom id", func() {
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
		})

		g.It("deletes existent user", func() {
			id := "d9dfc0c3-2bc4-4166-ba86-c7cc2818d554"
			_, err := usersCollection.InsertOne(ctx, bson.M{"name": "Alan", "active": true, "user_id": id})
			assert.Nil(t, err)

			err = userManager.Delete(id)

			assert.Nil(t, err)
			var saved models.User
			err = usersCollection.FindOne(ctx, bson.M{"user_id": id}).Decode(&saved)
			assert.Equal(t, "mongo: no documents in result", err.Error())
			assert.Equal(t, models.User{Name: "", Active: false, UserID: ""}, saved)
		})

		g.It("deletes non-existent user", func() {
			id := "non-existent"

			err := userManager.Delete(id)

			assert.Nil(t, err)
		})

		g.It("finds all users", func() {
			_, err := usersCollection.InsertMany(ctx, []interface{}{
				bson.M{"name": "Alan", "active": true, "user_id": "user-1"},
				bson.M{"name": "Bob", "active": false, "user_id": "user-2"},
			})
			assert.Nil(t, err)

			results, err := userManager.Find()

			assert.Nil(t, err)
			assert.Equal(t, []models.User{
				models.User{Name: "Alan", Active: true, UserID: "user-1"},
				models.User{Name: "Bob", Active: false, UserID: "user-2"},
			}, results)
		})

		g.It("finds empty list of users", func() {
			results, err := userManager.Find()

			assert.Nil(t, err)
			assert.Equal(t, []models.User{}, results)
		})

		g.It("filters by active", func() {
			_, err := usersCollection.InsertMany(ctx, []interface{}{
				bson.M{"name": "Alan", "active": true, "user_id": "user-1"},
				bson.M{"name": "Bob", "active": false, "user_id": "user-2"},
				bson.M{"name": "Cindy", "active": true, "user_id": "user-3"},
				bson.M{"name": "Dave", "active": false, "user_id": "user-4"},
			})
			assert.Nil(t, err)

			results, err := userManager.Find(map[string]interface{}{"active": true})

			assert.Nil(t, err)
			assert.Equal(t, []models.User{
				models.User{Name: "Alan", Active: true, UserID: "user-1"},
				models.User{Name: "Cindy", Active: true, UserID: "user-3"},
			}, results)
		})
	})
}
