package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/oprogramador/user-service-golang/controllers"
)

func HandleRequests(userManager UserManager) *gin.Engine {
	router := gin.Default()
	router.GET("/users", controllers.ListUsers(userManager))
	router.POST("/user", controllers.CreateUser(userManager))
	router.DELETE("/user/:id", controllers.DeleteUser(userManager))

	return router
}
