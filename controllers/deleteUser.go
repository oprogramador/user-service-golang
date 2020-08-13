package controllers

import (
	"github.com/gin-gonic/gin"
)

func DeleteUser(userManager UserManager) func(ginContext *gin.Context) {
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
