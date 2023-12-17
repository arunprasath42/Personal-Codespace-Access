package route

import (
	"gitworkflow-microservice/src/controllers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupRoutes(router *gin.Engine) {

	api := router.Group("/api/v1")

	api.POST("/cloneRepository", controllers.CloneRepository)
	api.POST("/createBranch", controllers.CreateBranch)
	api.POST("/commitChanges", controllers.CommitChanges)
	api.POST("/viewHistory", controllers.ViewHistory)
	api.POST("/showDiff", controllers.ShowDiff)

	router.Run(viper.GetString("server.port"))
}
