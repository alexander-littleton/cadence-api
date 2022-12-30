package main

import (
	"fmt"
	"github.com/alexander-littleton/cadence-api/configs"
	"github.com/alexander-littleton/cadence-api/internal/controllers"
	"github.com/alexander-littleton/cadence-api/internal/repositories"
	"github.com/alexander-littleton/cadence-api/internal/routes"
	userService "github.com/alexander-littleton/cadence-api/internal/services/user_service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	//TODO: setup trusted proxies
	router := gin.Default()
	validate := validator.New()

	routes.RegisterUserRoutes(
		router,
		controllers.NewUserController(
			userService.NewUserService(
				repositories.NewUserRepository(
					configs.GetCollection(configs.DB, "users"),
				),
				validate,
			),
		),
	)
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
