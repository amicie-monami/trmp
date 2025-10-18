package main

import (
	"trmp/internal/database"
	"trmp/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitSQLiteDB()
	if err != nil {
		panic(err)
	}

	serv := gin.Default()
	routes.SetupRoutes(serv, db)
	serv.Run()
}
