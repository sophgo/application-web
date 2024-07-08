package system

import (
	v1 "application-web/api/v1"

	"github.com/gin-gonic/gin"
)

type FileRouter struct{}

func (s *FileRouter) InitFileRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	fileRouter := Router.Group("file")
	fileApi := v1.ApiGroupApp.SystemApiGroup.FileApi
	{
		fileRouter.POST("upload", fileApi.Upload)
	}
	return fileRouter
}
