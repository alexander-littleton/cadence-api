package main

import (
	"github.com/gin-gonic/gin"
	"internal/controllers"
	"internal/repositories"
	"internal/routes"
)

func main() {
	router := gin.Default()
	//TODO: was there a better pattern for this?
	routes.RegisterUserRoutes(router, controllers.NewUserController(repositories.NewUserRepository()))
}
