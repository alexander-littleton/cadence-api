package routes

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser() gin.HandlerFunc
	GetUser() gin.HandlerFunc
	//EditAUser() gin.HandlerFunc
	//DeleteAUser() gin.HandlerFunc
}

func RegisterUserRoutes(router *gin.Engine, userController UserController) {
	router.POST("/user", userController.CreateUser())
	router.GET("/user/:userId", userController.GetUser())
	//router.PUT("/user/:userId", userController.EditAUser())
	//router.DELETE("/user/:userId", userController.DeleteAUser())
	//router.GET("/users", controllers.GetAllUsers())
}
