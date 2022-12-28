package main

import (
	"github.com/gin-gonic/gin"
	"internal/configs"
	"internal/controllers"
	"internal/repositories"
	"internal/routes"
	userService "internal/services"
)

func main() {
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
}
