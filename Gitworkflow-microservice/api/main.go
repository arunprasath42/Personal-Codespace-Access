package main

import (
	"gitworkflow-microservice/route"

	"gitworkflow-microservice/utils/middleware"

	config "gitworkflow-microservice/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	router := gin.Default()
	router.Use(middleware.TracingMiddleware())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
	route.SetupRoutes(router)
}
