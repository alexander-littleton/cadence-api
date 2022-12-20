package routes

import (
	"github.com/gin-gonic/gin"
	"internal/controllers"
)

//TODO: add interfaces for router and controller
func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	router.POST("/user", userController.CreateUser())
	//router.GET("/user/:userId", controllers.GetUser())
	router.PUT("/user/:userId", controllers.EditAUser())
	router.DELETE("/user/:userId", controllers.DeleteAUser())
	//router.GET("/users", controllers.GetAllUsers())
}
