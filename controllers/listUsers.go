package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/datamanagerinterfaces"
	"log"
)

func ListUsers(userManager datamanagerinterfaces.UserManager) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		users, err := userManager.Find()
		if err != nil {
			log.Println(err)
			ginContext.String(500, "")
		}
		ginContext.JSON(200, users)
	}
}
