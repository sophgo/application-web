package system

import (
	v1 "application-web/api/v1"
	"application-web/middleware"

	"github.com/gin-gonic/gin"
)

type FaceTaskRouter struct {
}

func (s *FaceTaskRouter) InitFaceTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("api/face")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.FaceTaskApi
	{
		router.POST("add", api.CreateTask)
		router.POST("delete", api.RemoveTask)
		router.POST("start", api.StartTask)
		router.POST("stop", api.StopTask)
		router.POST("list", api.ListTask)
		router.POST("compare", api.Compare)
		router.POST("feature", api.GetFeature)
		router.POST("gallery/create", api.CreateGallery)

	}

	return router
}

func (s *FaceTaskRouter) InitFaceUploadRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	router := Router.Group("api/face")
	router.Use(middleware.BlockerMiddleware())
	api := v1.ApiGroupApp.SystemApiGroup.FaceTaskApi
	{

		router.POST("upload", api.ReceiveResult)
		router.POST("upload/search", api.ReceiveSearch)
		router.POST("search/add", api.CreateSearchTask)
		router.POST("library/add", api.CreatePersonInfo)
		router.POST("library/delete", api.RemovePersonInfo)
		router.GET("library/list", api.ListPersonInfo)
		router.POST("compareTest", api.CompareTest)

	}

	return router
}
