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
			if err != nil {
				log.Println(err)
				ginContext.String(400, "'active' param in the query string should be of type bool")
				return
			}
			userManagerQuery = map[string]interface{}{"active": active}
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
			return
		}
		ginContext.JSON(200, users)
	}
}
