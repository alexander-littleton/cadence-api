package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"internal/configs"
	"internal/controllers"
	"internal/repositories"
	"internal/routes"
	userService "internal/services"
)

func main() {
	//TODO: setup trusted proxies
	router := gin.Default()

	routes.RegisterUserRoutes(
		router,
		controllers.NewUserController(
			userService.NewUserService(
				repositories.NewUserRepository(
					configs.GetCollection(configs.DB, "users"),
				),
			),
		),
	)
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
