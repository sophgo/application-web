package system

import (
	v1 "application-web/api/v1"
	"application-web/middleware"

	"github.com/gin-gonic/gin"
)

type ReceiveAlarmRouter struct{}

func (s *ReceiveAlarmRouter) InitReceiveAlarmRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.ReceiveAlarmApi
	{
		router.POST("upload", api.AlarmRev)
	}

	return router
}
