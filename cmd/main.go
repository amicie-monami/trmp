package cmd

import (
	"github.com/gin-gonic/gin"
)

func main() {
	serv := gin.Default()
	serv.Run()
}
