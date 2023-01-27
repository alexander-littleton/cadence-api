package main

import (
	"fmt"
	"github.com/alexander-littleton/cadence-api/configs"
	userService "github.com/alexander-littleton/cadence-api/pkg/user"
	"github.com/alexander-littleton/cadence-api/pkg/user/api"
	"github.com/alexander-littleton/cadence-api/pkg/user/repositories"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	//TODO: setup trusted proxies
	router := gin.Default()
	validate := validator.New()
	userController := api.New(
		userService.New(
			repositories.NewUserRepository(
				configs.GetCollection(configs.DB, "users"),
			),
			validate,
		),
	)
	userController.RegisterRoutes(router)
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
