package routes

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	GetUserByEmail(ctx *gin.Context)
	//EditAUser() gin.HandlerFunc
	//DeleteAUser() gin.HandlerFunc
}

func RegisterUserRoutes(router *gin.Engine, userController UserController) {
	router.POST("/user", userController.CreateUser)
	router.GET("/user/:userId", userController.GetUserById)
	router.GET("/user/:email", userController.GetUserByEmail)
	//router.PUT("/user/:userId", userController.EditAUser())
	//router.DELETE("/user/:userId", userController.DeleteAUser())
	//router.GET("/users", controllers.GetAllUsers())
}
