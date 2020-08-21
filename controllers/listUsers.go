package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/datamanagerinterfaces"
	"github.com/oprogramador/user-service-golang/models"
	"log"
	"strconv"
)

func ListUsers(userManager datamanagerinterfaces.UserManager) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		query := ginContext.Request.URL.Query()
		queryActive := query["active"]
		userManagerQuery := map[string]interface{}(nil)
		if len(queryActive) > 0 {
			active, err := strconv.ParseBool(queryActive[0])
			userManagerQuery = map[string]interface{}{"active": active}
			if err != nil {
				log.Println(err)
				ginContext.String(500, "")
			}
		}
		var err error
		var users []models.User
		if userManagerQuery == nil {
			users, err = userManager.Find()
		} else {
			users, err = userManager.Find(userManagerQuery)
		}
		if err != nil {
			log.Println(err)
			ginContext.String(500, "")
		}
		ginContext.JSON(200, users)
	}
}
