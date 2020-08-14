package models

type User struct {
	UserID string `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Active bool   `json:"active" bson:"active"`
}
