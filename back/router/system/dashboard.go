package system

import (
	v1 "application-web/api/v1"
	"application-web/middleware"

	"github.com/gin-gonic/gin"
)

type DashBoardRouter struct{}

func (s *DashBoardRouter) InitBoardRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	router := Router.Group("api/dashboard")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.DashBoardApi
	{
		router.GET("info", api.LoadDashboardBaseInfo)

	}

	return router
}
