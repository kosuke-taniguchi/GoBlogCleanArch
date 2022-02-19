package main

import (
	"github.com/gin-gonic/gin"

	"GoBlogClean/article/handler"
	"GoBlogClean/injector"
)

func main() {
	articleHandler := injector.InjectArticleHandler()
	engine := gin.Default()
	handler.InitRouting(engine, &articleHandler)
	engine.Run(":8080")
}