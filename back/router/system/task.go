package system

import (
	v1 "application-web/api/v1"
	"application-web/middleware"

	"github.com/gin-gonic/gin"
)

type TaskRouter struct{}

func (s *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api/task")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.TaskApi
	{
		router.GET("abilities", api.AbilitiesList)

		router.POST("add", api.AddTask)
		router.POST("modify", api.ModTask)
		router.POST("delete", api.DeleteTask)
		router.POST("start", api.StartTask)
		router.POST("stop", api.StopTask)
		router.POST("list", api.List)

		router.POST("image", api.UpdateImage)

	}

	return router
}
func (s *TaskRouter) InitImageRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api/task")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.TaskApi
	{
		router.GET("image", api.GetImage)

	}

	return router
}
func (s *TaskRouter) InitConfigRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api/config")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.TaskApi
	{
		router.POST("get", api.GetTaskConfig)
		router.POST("mod", api.ModTaskConfig)

	}

	return router
}
