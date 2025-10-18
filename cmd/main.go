package main

import (
	"trmp/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	serv := gin.Default()
	routes.SetupRoutes(serv)
	serv.Run()
}
