package main

import (
	"fmt"

	"github.com/DeniesKresna/amarthatest/config"
	"github.com/DeniesKresna/amarthatest/one"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.Config{}

	err := appConfig.InitDatabase()
	if err != nil {
		fmt.Println(err)
	}

	r := gin.Default()
	one.InitQuestion(r, &appConfig)

	r.Run("localhost:8898")
}
