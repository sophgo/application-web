package system

import (
	v1 "application-web/api/v1"
	"application-web/middleware"

	"github.com/gin-gonic/gin"
)

type QueryRouter struct{}

func (s *QueryRouter) InitQueryRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	// router := Router.Group("algorithm/alarm", middleware.TimeoutMiddleware(global.TimeOut))
	router := Router.Group("api/alarm")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.QueryApi
	{
		router.POST("list", api.GetRecord)
		router.GET("info", api.GetQueryInfo)
		router.POST("modSize", api.ModSize)
		router.POST("delete", api.DeleteAlarms)
	}

	return router
}

func (s *QueryRouter) InitGetImageRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	// router := Router.Group("algorithm/alarm", middleware.TimeoutMiddleware(global.TimeOut))
	router := Router.Group("api/alarm")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.QueryApi
	{
		router.GET("image", api.GetImage)

	}

	return router
}
