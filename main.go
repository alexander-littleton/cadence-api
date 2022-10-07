package main

import (
	"github.com/gin-gonic/gin"
	"internal/routes"
)

func main() {
	router := gin.Default()
	routes.UserRoute(router)
}
