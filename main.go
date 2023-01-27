package main

import (
	"fmt"
	"github.com/alexander-littleton/cadence-api/configs"
	userService "github.com/alexander-littleton/cadence-api/pkg/user"
	userApi "github.com/alexander-littleton/cadence-api/pkg/user/api"
	userRepo "github.com/alexander-littleton/cadence-api/pkg/user/repositories/mongo"
	"github.com/gin-gonic/gin"
)

func main() {
	//TODO: setup trusted proxies
	router := gin.Default()
	userController := userApi.New(
		userService.New(
			userRepo.NewUserRepository(
				configs.GetCollection(configs.DB, "users"),
			),
		),
	)
	userController.RegisterRoutes(router)
	err := router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
