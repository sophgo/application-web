package initialization

import (
	"application-web/dist"
	"application-web/logger"
	"application-web/middleware"
	"application-web/router"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	gin.SetMode(gin.ReleaseMode) // 设置Gin的模式为release
	Router := gin.New()
	Router.Use(gin.Recovery())

	systemRouter := router.RouterGroupApp.System

	cors_config := SetCors()
	Router.Use(middleware.BlockerMiddleware())
	Router.Use(cors.New(cors_config))
	// 设置404处理器
	Router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	PublicGroup := Router.Group("")
	{
		setWebStatic(PublicGroup)
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		systemRouter.InitGetImageRouter(PublicGroup)
		systemRouter.InitImageRouter(PublicGroup)
		systemRouter.InitFaceImageRouter(PublicGroup)
		systemRouter.InitFaceUploadRouter(PublicGroup)
		systemRouter.InitFileRouter(PublicGroup)
		systemRouter.InitReceiveAlarmRouter(PublicGroup)
	}

	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.AuthMiddleware())
	{
		systemRouter.InitTaskRouter(PrivateGroup)
		systemRouter.InitFaceTaskRouter(PrivateGroup)
		// systemRouter.InitFileRouter(PrivateGroup) todo
		systemRouter.InitQueryRouter(PrivateGroup)
		systemRouter.InitConfigRouter(PrivateGroup)
		systemRouter.InitFaceAlarmRouter(PrivateGroup)
		systemRouter.InitBoardRouter(PrivateGroup)
	}
	logger.Info("Router Init Ok")
	return Router
}

func setWebStatic(rootRouter *gin.RouterGroup) {
	rootRouter.Use(func(c *gin.Context) {
		c.Next()
	})

	rootRouter.GET("/assets/*filepath", func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d", 3600))
		staticServer := http.FileServer(http.FS(dist.Assets))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	rootRouter.GET("/", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(dist.IndexHtml))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	rootRouter.GET("/admin.gif", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(dist.AdminGif))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	rootRouter.GET("/logo.png", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(dist.LogoPng))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	rootRouter.GET("/favicon.ico", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(dist.Favicon))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
}

func SetCors() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	return config
}
