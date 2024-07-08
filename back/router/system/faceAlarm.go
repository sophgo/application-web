package system

import (
	v1 "application-web/api/v1"
	"application-web/middleware"

	"github.com/gin-gonic/gin"
)

type FaceAlarmRouter struct{}

func (s *FaceAlarmRouter) InitFaceAlarmRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api/face/alarm")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.FaceAlarmApi
	{
		router.POST("list", api.List)
		router.GET("info", api.GetQueryInfo)
		router.POST("modSize", api.ModSize)
		router.POST("delete", api.DeleteAlarms)
		router.POST("search/res/list", api.ListSearchResult)
	}

	return router
}
func (s *FaceAlarmRouter) InitFaceImageRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api/face/alarm")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.FaceAlarmApi
	{
		router.GET("image", api.GetImage)

	}

	return router
}
