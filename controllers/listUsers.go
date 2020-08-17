package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func ListUsers(userManager UserManager) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		users, err := userManager.Find()
		if err != nil {
			log.Println(err)
			ginContext.String(500, "")
		}
		ginContext.JSON(200, users)
	}
}
