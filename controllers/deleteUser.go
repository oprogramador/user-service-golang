package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/datamanagerinterfaces"
)

func DeleteUser(userManager datamanagerinterfaces.UserManager) func(ginContext *gin.Context) {
	return func(ginContext *gin.Context) {
		id := ginContext.Param("id")

		err := userManager.Delete(id)
		if err != nil {
			ginContext.String(400, "invalid id")
			return
		}

		ginContext.String(204, "")
	}
}
