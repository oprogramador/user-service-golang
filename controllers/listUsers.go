package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/datamanagerinterfaces"
	"log"
	"strconv"
)

func ListUsers(userManager datamanagerinterfaces.UserManager) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		query := ginContext.Request.URL.Query()
		active, err := strconv.ParseBool(query["active"][0])
		if err != nil {
			log.Println(err)
			ginContext.String(500, "")
		}
		users, err := userManager.Find(map[string]interface{}{"active": active})
		if err != nil {
			log.Println(err)
			ginContext.String(500, "")
		}
		ginContext.JSON(200, users)
	}
}
