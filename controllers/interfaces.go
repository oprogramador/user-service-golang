package controllers

import (
	"github.com/oprogramador/user-service-golang/models"
)

type UserManager interface {
	Save(user *models.User) error
	Delete(id string) error
	FindAll() ([]models.User, error)
}
