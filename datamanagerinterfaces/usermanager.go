package datamanagerinterfaces

import (
	"github.com/oprogramador/user-service-golang/models"
)

type UserManager interface {
	Save(user *models.User) error
	Delete(id string) error
	Find(...map[string]interface{}) ([]models.User, error)
}
